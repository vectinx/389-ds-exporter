package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	slogmulti "github.com/samber/slog-multi"

	"389-ds-exporter/src/cmd"
	"389-ds-exporter/src/collectors"
	"389-ds-exporter/src/config"
	"389-ds-exporter/src/metrics"
	"389-ds-exporter/src/network"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// This variables is filled via ldflags at build time
	Version    = "dev"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)

// appResources struct contains pointers to resources that must be closed when the program terminates.
// Resources must be added to the structure as they are initialized
type appResources struct {
	LogFile    *os.File
	ConnPool   *network.LdapConnectionPool
	HttpServer *http.Server
}

func (r *appResources) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down gracefully...")

	var errs []error
	if r.HttpServer != nil {
		slog.Debug("Stopping HTTP server ...")
		r.HttpServer.SetKeepAlivesEnabled(false)
		if err := r.HttpServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("HTTP server Shutdown failed: %v", err))
		}
		slog.Debug("HTTP server stopped")
	}

	if r.ConnPool != nil {
		slog.Debug("Closing LDAP connection pool ...")
		if err := r.ConnPool.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("error closing ldap pool: %v", err))
		}
		slog.Debug("LDAP connection closed")
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
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

func reopenLogFile(cfg *config.ExporterConfiguration, resources *appResources) error {
	if resources.LogFile != nil {
		_ = resources.LogFile.Close()
	}

	logFile, err := os.OpenFile(cfg.Logging.GetFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	newLogger, _, err := setupLogger(cfg)
	if err != nil {
		logFile.Close()
		return fmt.Errorf("failed to set up new logger: %w", err)
	}

	slog.SetDefault(newLogger)
	resources.LogFile = logFile

	return nil
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

func setupPrometheusMetrics(cfg *config.ExporterConfiguration, connPool *network.LdapConnectionPool) *prometheus.Registry {
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

func run() int {
	var (
		applicationResources = appResources{}
		args                 = cmd.ParseCmdArguments(fmt.Sprintf("Version: %s\nCommit: %s\nBuild time: %s", Version, CommitHash, BuildTime))

		signalCh    = make(chan os.Signal, 1)
		serverErrCh = make(chan error)
	)

	cfg, err := readConfig(args.ConfigFile)
	if err != nil {
		slog.Error("Error loading config", "err", err)
		return 1
	}

	if args.IsConfigCheck {
		fmt.Print(cfg.String())
		return 0
	}

	slog.Info("Configuration read successfuly")

	defer func() {
		shutdownContext, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Global.GetShutdownTimeout()*uint(time.Second)))
		defer cancel()
		if err := applicationResources.Shutdown(shutdownContext); err != nil {
			slog.Error("Shutdown error", "err", err)
		}
		if applicationResources.LogFile != nil {
			// We close the log file last, because we expect logs to be written to it until the very end
			if err := applicationResources.LogFile.Close(); err != nil {
				slog.Error("Error closing log file", "err", err)
			}
		}
	}()

	logger, logFile, err := setupLogger(cfg)
	applicationResources.LogFile = logFile

	if err != nil {
		slog.Error("Error initializing logging", "err", err)
		return 1
	}
	slog.SetDefault(logger)

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	slog.Info("Starting 389-ds-exporter", "version", Version, "commit", CommitHash, "build_time", BuildTime)
	slog.Info("LDAP server info",
		"url", cfg.LDAP.ServerURL,
		"bind_dn", cfg.LDAP.BindDN,
		"backend", cfg.Global.BackendImplement,
	)

	ldapConnPoolConfig := network.LdapConnectionPoolConfig{
		ServerURL:              cfg.LDAP.ServerURL,
		BindDN:                 cfg.LDAP.BindDN,
		BindPw:                 cfg.LDAP.BindPw,
		MaxConnections:         cfg.LDAP.ConnectionPool.GetConnectionsLimit(),
		DialTimeout:            time.Duration(cfg.LDAP.ConnectionPool.GetDialTimeout()) * time.Second,
		RetryCount:             cfg.LDAP.ConnectionPool.GetRetryCount(),
		RetryDelay:             time.Duration(cfg.LDAP.ConnectionPool.GetRetryDelay()) * time.Second,
		ConnectionAliveTimeout: time.Duration(cfg.LDAP.ConnectionPool.GetConnectionAliveTimeout()) * time.Second,
	}

	ldapConnPool := network.NewLdapConnectionPool(ldapConnPoolConfig)
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
	ln, _ := net.Listen("tcp", cfg.HTTP.GetListenAddress())
	timeoutListener := network.NewTimeoutListener(ln, time.Duration(cfg.HTTP.GetInitialReadTimeout()*uint(time.Second)))

	go func() {
		slog.Info(fmt.Sprintf("Starting HTTP server at %s", cfg.HTTP.GetListenAddress()))
		if err := server.Serve(timeoutListener); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	running := true

	for running {
		select {
		case signal := <-signalCh:
			switch signal {
			case syscall.SIGINT:
				slog.Info("SIGINT signal received")
				running = false
			case syscall.SIGTERM:
				slog.Info("SIGTERM signal received")
				running = false
			case syscall.SIGHUP:
				slog.Info("SIGHUP signal received - reopening log file")
				if reopenLogFile(cfg, &applicationResources) != nil {
					slog.Error("Error reopening log file")
					running = false
				}
			}
		case err := <-serverErrCh:
			slog.Error(fmt.Sprintf("HTTP server failed with error: %v", err))
			running = false
		}
	}
	return 0
}

func main() {
	os.Exit(run())
}
