package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	slogmulti "github.com/samber/slog-multi"

	"389-ds-exporter/src/backends"
	"389-ds-exporter/src/collectors"
	"389-ds-exporter/src/config"
	"389-ds-exporter/src/metrics"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// This variables is filled via ldflags at build time
	Version   = "dev"
	BuildTime = "unknown"
)

// appResources struct contains pointers to resources that must be closed when the program terminates.
// Resources must be added to the structure as they are initialized
type appResources struct {
	LogFile    *os.File
	ConnPool   *backends.LdapConnectionPool
	HttpServer *http.Server
}

func (r *appResources) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down gracefully...")

	if r.HttpServer != nil {
		slog.Debug("Stopping HTTP server ...")
		if err := r.HttpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("HTTP server Shutdown failed: %v", err)
		}
		slog.Debug("HTTP server stopped")
	}

	if r.ConnPool != nil {
		slog.Debug("Closing LDAP connection pool ...")
		if err := r.ConnPool.Close(ctx); err != nil {
			return fmt.Errorf("error closing ldap pool: %v", err)
		}
		slog.Debug("LDAP connection closed")
	}

	if r.LogFile != nil {
		if err := r.LogFile.Close(); err != nil {
			return fmt.Errorf("error closing log file: %v", err)
		}
	}
	slog.Info("All resources shut down successfully")
	return nil
}

// defaultHttpResponse function generates a standard HTML response for the exporter
func defaultHttpResponse(metricsPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, `
		<html>
			<head>
				<title>389-ds-exporter</title>
			</head>
			<body>
				<p>Metrics are <a href="%s">here</a></p>
			</body>
		</html>`, html.EscapeString(metricsPath))
		if err != nil {
			log.Printf("Error writing HTTP answer: %s", err)
		}
	}
}

func buildLogHandler(format string, w io.Writer, level slog.Level) slog.Handler {
	if format == "text" {
		return slog.NewTextHandler(w, &slog.HandlerOptions{Level: level})
	} else if format == "json" {
		return slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})
	}
	return nil
}

func setupLogger(cfg *config.ExporterConfiguration) (*slog.Logger, *os.File, error) {
	var logLevel slog.Level
	handlers := []slog.Handler{}
	var logFile *os.File

	strLogLevel := cfg.Logging.GetLevel()
	levelMap := map[string]slog.Level{
		"DEBUG":   slog.LevelDebug,
		"INFO":    slog.LevelInfo,
		"WARNING": slog.LevelWarn,
		"ERROR":   slog.LevelError,
	}

	logLevel, ok := levelMap[strLogLevel]
	if !ok {
		return nil, nil, fmt.Errorf("unknown logging level: '%s'", strLogLevel)
	}

	if cfg.Logging.GetHandler() == "stdout" || cfg.Logging.GetHandler() == "both" {
		handlers = append(handlers, buildLogHandler(cfg.Logging.GetStdoutFormat(), os.Stdout, logLevel))
	}
	if cfg.Logging.GetHandler() == "file" || cfg.Logging.GetHandler() == "both" {
		var err error
		logFile, err = os.OpenFile(cfg.Logging.GetFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("error opening log file: %w", err)
		}
		handlers = append(handlers, buildLogHandler(cfg.Logging.GetFileFormat(), logFile, logLevel))
	}

	if len(handlers) == 0 {
		return nil, nil, errors.New("unable to create logger - logging handlers not specified")
	}

	logger := slog.New(slogmulti.Fanout(handlers...))
	return logger, logFile, nil
}

func readConfig(configFilePath string) (*config.ExporterConfiguration, error) {
	configuration, err := config.ReadConfig(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration: %w", err)
	}

	err = configuration.Validate()
	if err != nil {
		return nil, fmt.Errorf("incorrect configuration provided: %w", err)
	}
	return &configuration, nil
}

func setupPrometheusMetrics(cfg *config.ExporterConfiguration, connPool *backends.LdapConnectionPool) *prometheus.Registry {
	dsMetricsRegistry := prometheus.NewRegistry()

	dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
		"ds_exporter",
		connPool,
		"cn=monitor",
		metrics.GetLdapServerMetrics(),
		prometheus.Labels{},
	),
	)

	dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
		"ds_exporter",
		connPool,
		"cn=snmp,cn=monitor",
		metrics.GetLdapServerSnmpMetrics(),
		prometheus.Labels{},
	),
	)

	for _, backend := range cfg.Global.Backends {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			connPool,
			"cn=monitor,cn="+backend+",cn=ldbm database,cn=plugins,cn=config",
			metrics.GetLdapBackendCaches(),
			prometheus.Labels{"database": backend},
		),
		)
	}

	/*
		Since 389-ds has a different set of monitoring metrics for different backends (Berkley DB and LMDB),
		at the initialization stage we select the metrics that correspond to the selected backend
	*/
	if cfg.Global.BackendImplement == config.BackendBDB {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			connPool,
			"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
			metrics.GetLdapBDBServerCacheMetrics(),
			prometheus.Labels{},
		),
		)
	} else if cfg.Global.BackendImplement == config.BackendMDB {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			connPool,
			"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
			metrics.GetLdapMDBServerCacheMetrics(),
			prometheus.Labels{},
		),
		)
	}
	return dsMetricsRegistry
}

func main() {
	var (
		configFilePath = kingpin.Flag("config", "Path to configuration file").
				Default("config.yml").
				String()
		checkConfig             = kingpin.Flag("check-config", "Check current configuration and print it to stdout").Bool()
		shutdownContext, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		applicationResources    = appResources{}
		signalCh                = make(chan os.Signal, 1)
		serverErrCh             = make(chan error)
	)

	defer cancel()
	defer func() {
		if err := applicationResources.Shutdown(shutdownContext); err != nil {
			slog.Error("Shutdown error", "err", err)
		}
	}()

	kingpin.Version(fmt.Sprintf("Version: %s\nBuild time: %s", Version, BuildTime))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	cfg, err := readConfig(*configFilePath)
	if err != nil {
		slog.Error("Error loading config: %v", "err", err)
		os.Exit(1)
	}

	if *checkConfig {
		fmt.Print(cfg.String())
		return
	}

	logger, logFile, err := setupLogger(cfg)

	applicationResources.LogFile = logFile

	if err != nil {
		slog.Error("Error initializing logging", "err", err)
		return
	}
	slog.SetDefault(logger)

	slog.Info("Starting 389-ds-exporter", "version", Version, "build_time", BuildTime)
	slog.Info("Configuration read successfuly")
	slog.Info("LDAP server info",
		"url", cfg.LDAP.ServerURL,
		"bind_dn", cfg.LDAP.BindDN,
		"backend", cfg.Global.BackendImplement,
	)

	ldapConnPoolConfig := backends.LdapConnectionPoolConfig{
		ServerURL:              cfg.LDAP.ServerURL,
		BindDN:                 cfg.LDAP.BindDN,
		BindPw:                 cfg.LDAP.BindPw,
		MaxConnections:         cfg.LDAP.ConnectionPool.GetConnectionsLimit(),
		DialTimeout:            time.Duration(cfg.LDAP.ConnectionPool.GetDialTimeout()) * time.Second,
		RetryCount:             cfg.LDAP.ConnectionPool.GetRetryCount(),
		RetryDelay:             time.Duration(cfg.LDAP.ConnectionPool.GetRetryDelay()) * time.Second,
		ConnectionAliveTimeout: time.Duration(cfg.LDAP.ConnectionPool.GetConnectionAliveTimeout()) * time.Second,
	}

	ldapConnPool := backends.NewLdapConnectionPool(ldapConnPoolConfig)

	applicationResources.ConnPool = ldapConnPool

	dsMetricsRegistry := setupPrometheusMetrics(cfg, ldapConnPool)

	http.Handle(cfg.HTTP.GetMetricsPath(), promhttp.HandlerFor(dsMetricsRegistry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", defaultHttpResponse(cfg.HTTP.GetMetricsPath()))

	server := &http.Server{
		Addr:         cfg.HTTP.GetListenAddress(),
		Handler:      http.DefaultServeMux,
		ReadTimeout:  time.Duration(cfg.HTTP.GetReadTimeout()) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTP.GetWriteTimeout()) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTP.GetIdleTimeout()) * time.Second,
	}

	applicationResources.HttpServer = server

	go func() {
		slog.Info(fmt.Sprintf("Starting HTTP server at %s", cfg.HTTP.GetListenAddress()))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case signal := <-signalCh:
		switch signal {
		case syscall.SIGINT:
			slog.Info("SIGINT signal received")
		case syscall.SIGTERM:
			slog.Info("SIGTERM signal received")
		}
	case err := <-serverErrCh:
		slog.Error(fmt.Sprintf("HTTP server failed with error: %v", err))
	}
}
