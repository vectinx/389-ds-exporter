package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"389-ds-exporter/src/cmd"
	"389-ds-exporter/src/config"
	"389-ds-exporter/src/connections"
	"389-ds-exporter/src/metrics"
	"389-ds-exporter/src/utils"
)

var (
	// This variables is filled via ldflags at build time.
	Version    = "dev"     //nolint:gochecknoglobals
	BuildTime  = "unknown" //nolint:gochecknoglobals
	CommitHash = "unknown" //nolint:gochecknoglobals
)

const (
	LdapConnectionPoolTimeout time.Duration = 3 * time.Second
)

// appResources struct contains pointers to resources that must be closed when the program terminates.
// Resources must be added to the structure as they are initialized.
type appResources struct {
	LogFile    *os.File
	ConnPool   *connections.LdapConnectionPool
	HttpServer *http.Server
}

func (r *appResources) shutdown(ctx context.Context) error {
	slog.Info("Shutting down gracefully...")

	var errs []error
	if r.HttpServer != nil {
		slog.Debug("Stopping HTTP server ...")
		r.HttpServer.SetKeepAlivesEnabled(false)
		err := r.HttpServer.Shutdown(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("HTTP server Close failed: %w", err))
		}
		slog.Debug("HTTP server stopped", "err", err)
	}

	if r.ConnPool != nil {
		slog.Debug("Closing LDAP connection pool ...")
		err := r.ConnPool.Close(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("error closing ldap pool: %w", err))
		}
		slog.Debug("LDAP connection pool closed", "err", err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	slog.Info("All resources shut down successfully")

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

func run() int {
	var (
		applicationResources = appResources{}
		startTime            = time.Now()
		args                 = cmd.ParseCmdArguments(
			fmt.Sprintf(
				"Version: %s\nCommit: %s\nBuild time: %s",
				Version,
				CommitHash,
				BuildTime,
			),
		)
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

	slog.Info("Configuration read successfully")

	defer func() {
		shutdownContext, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.Global.GetShutdownTimeout())*time.Second,
		)

		defer cancel()
		err := applicationResources.shutdown(shutdownContext)
		if err != nil {
			slog.Error("Shutdown error", "err", err)
		}
		if applicationResources.LogFile != nil {
			// We close the log file last, because we expect logs to be written to it until the very end
			err := applicationResources.LogFile.Close()
			if err != nil {
				slog.Error("Error closing log file", "err", err)
			}
		}
	}()

	logger, logFile, err := utils.SetupLogger(cfg)
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

	ldapConnPoolConfig := connections.LdapConnectionPoolConfig{
		ServerURL:      cfg.LDAP.ServerURL,
		BindDN:         cfg.LDAP.BindDN,
		BindPw:         cfg.LDAP.BindPw,
		MaxConnections: cfg.LDAP.GetPoolConnLimit(),
		ConnFactory:    connections.RealConnectionDialUrl,
	}

	applicationResources.ConnPool = connections.NewLdapConnectionPool(ldapConnPoolConfig)

	dsMetricsRegistry := metrics.SetupPrometheusMetrics(cfg, applicationResources.ConnPool, LdapConnectionPoolTimeout)

	// Create HTTP server
	applicationResources.HttpServer = &http.Server{
		Addr:         cfg.HTTP.GetListenAddress(),
		Handler:      http.DefaultServeMux,
		ReadTimeout:  time.Duration(cfg.HTTP.GetReadTimeout()) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTP.GetWriteTimeout()) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTP.GetIdleTimeout()) * time.Second,
	}

	// Create HTTP Listener with timeouts
	listener, err := net.Listen("tcp", cfg.HTTP.GetListenAddress())
	if err != nil {
		slog.Error("Failed to start TCP listener", "err", err)
		return 1
	}
	timeoutListener := connections.NewTimeoutListener(
		listener,
		time.Duration(cfg.HTTP.GetInitialReadTimeout())*time.Second,
	)

	// Register HTTP endpoinnts
	http.Handle(cfg.HTTP.GetMetricsPath(), promhttp.HandlerFor(dsMetricsRegistry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", utils.DefaultHttpResponse(cfg.HTTP.GetMetricsPath()))
	http.HandleFunc("/health", utils.HealthHttpResponse(applicationResources.ConnPool, startTime))
	http.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Start HTTP server
	go func() {
		slog.Info("Starting HTTP server at " + cfg.HTTP.GetListenAddress())
		err := applicationResources.HttpServer.Serve(timeoutListener)
		if err != nil && err != http.ErrServerClosed {
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
				slog.Info("SIGHUP signal received")
				if applicationResources.LogFile != nil {
					slog.Info("Reopening log file")
					newLogFile, err := utils.ReopenLogFile(cfg, applicationResources.LogFile)
					if err != nil {
						slog.Error("Error reopening log file")
						running = false
					}
					applicationResources.LogFile = newLogFile
					slog.Info("Log file reopenedd successfully")
				}
			}
		case err := <-serverErrCh:
			slog.Error(fmt.Sprintf("HTTP server failed with error: %v", err))
			running = false
		}
	}

	// Before return, deferred Shutdown function will be executed
	return 0
}

func main() {
	os.Exit(run())
}
