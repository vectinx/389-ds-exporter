package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"

	"389-ds-exporter/src/cmd"
	"389-ds-exporter/src/config"
	expldap "389-ds-exporter/src/ldap"
	"389-ds-exporter/src/metrics"
	"389-ds-exporter/src/utils"
)

// appResources struct contains pointers to resources that must be closed when the program terminates.
// Resources must be added to the structure as they are initialized.
type appResources struct {
	ConnPool   *expldap.LDAPPool
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
		err := r.ConnPool.Close()
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

func readConfig(configFilePath string) (*config.ExporterConfig, error) {
	configuration, err := config.ReadConfig(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration: %w", err)
	}

	err = configuration.Validate()
	if err != nil {
		return nil, fmt.Errorf("incorrect configuration provided: %w", err)
	}

	return configuration, nil
}

func run() int {
	var (
		applicationResources = appResources{}
		startTime            = time.Now()
		args                 = cmd.ParseCmdArguments()
		signalCh             = make(chan os.Signal, 1)
		serverErrCh          = make(chan error)
	)

	logger := promslog.New(args.PromslogConfig)
	slog.SetDefault(logger)

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
			time.Duration(cfg.ShutdownTimeout)*time.Second,
		)

		defer cancel()
		err := applicationResources.shutdown(shutdownContext)
		if err != nil {
			slog.Error("Shutdown error", "err", err)
		}
	}()

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	slog.Info(
		"Starting 389-ds-exporter",
		"version", version.Version,
		"commit", version.Revision,
		"build_time", version.BuildDate,
	)

	slog.Info("LDAP server info", "url", cfg.LDAPServerURL, "bind_dn", cfg.LDAPBindDN)

	ldapConnPoolConfig := expldap.LDAPPoolConfig{
		Auth: expldap.LDAPAuthConfig{
			URL:           cfg.LDAPServerURL,
			BindDN:        cfg.LDAPBindDN,
			BindPw:        cfg.LDAPBindPw,
			DialTimeout:   time.Duration(cfg.LDAPDialTimeout) * time.Second,
			TlsSkipVerify: cfg.LDAPTlsSkipVerify,
		},
		DialTimeout:    time.Duration(cfg.LDAPDialTimeout) * time.Second,
		MaxConnections: cfg.LDAPPoolConnLimit,
		MaxIdleTime:    time.Duration(cfg.LDAPPoolIdleTime) * time.Second,
		MaxLifeTime:    time.Duration(cfg.LDAPPoolLifeTime) * time.Second,
		ConnFactory:    expldap.RealConnectionDialUrl,
	}

	applicationResources.ConnPool = expldap.NewLDAPPool(ldapConnPoolConfig)

	dsMetricsRegistry := metrics.SetupPrometheusMetrics(
		cfg,
		applicationResources.ConnPool,
	)

	// Create HTTP server
	// #nosec G112: HTTP timeouts will be configured later using the exporter-toolkit
	applicationResources.HttpServer = &http.Server{}

	// Register HTTP endpoinnts
	http.Handle(args.MetricsPath, promhttp.HandlerFor(dsMetricsRegistry, promhttp.HandlerOpts{}))

	if args.MetricsPath != "/" {
		landingConfig := web.LandingConfig{
			Name:        "389-ds-exporter",
			Description: "Prometheus exporter for 389 Directory Server",
			Version:     version.Info(),
			HeaderColor: "#0c7982",
			Profiling:   "false",
			Links: []web.LandingLinks{
				{
					Address: args.MetricsPath,
					Text:    "Metrics",
				},
				{
					Address: "/health",
					Text:    "Health",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			slog.Error(err.Error())
			return 1
		}
		http.Handle("/", landingPage)
	}

	http.HandleFunc("/health", utils.HealthHttpResponse(
		applicationResources.ConnPool,
		startTime,
		time.Duration(cfg.LDAPPoolGetTimeout)*time.Second),
	)

	http.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	})

	// Start HTTP server
	go func() {
		err := web.ListenAndServe(applicationResources.HttpServer, args.ExporterToolkitFlags, logger)
		if !errors.Is(err, http.ErrServerClosed) {
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
