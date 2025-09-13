package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// defaultHttpResponse function generates a standard HTML response for the exporter
func defaultHttpResponse(metricsPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`
		<html>
			<head>
				<title>389-ds-exporter</title>
			</head>
			<body>
				<p>Metrics are <a href='` + metricsPath + `'>here</a></p>
			</body>
		</html>`))
		if err != nil {
			log.Printf("Error writing HTTP answer: %s", err)
		}
	}
}

func main() {
	var (
		configFilePath = kingpin.Flag("config", "Path to configuration file").
				Default("config.yml").
				String()
		showConfig = kingpin.Flag("check-config", "Check current configuration and print it to stdout").Bool()
	)

	kingpin.Version(fmt.Sprintf("Version: %s\nBuild time: %s", Version, BuildTime))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	configuration, err := config.ReadConfig(*configFilePath)
	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	err = configuration.Validate()
	if err != nil {
		log.Fatalf("Incorrect configuration provided: %v", err)
	}

	if *showConfig {
		fmt.Print(configuration.String())
		return
	}

	log.Printf("Configuration read successfuly")
	log.Printf("LDAP server URL: %v", configuration.LDAP.ServerURL)
	log.Printf("LDAP bind DN: %v", configuration.LDAP.BindDN)
	log.Printf("389-ds backend type: %v", configuration.Global.BackendImplement)

	dsMetricsRegistry := prometheus.NewRegistry()

	ldapConnPoolConfig := backends.LdapConnectionPoolConfig{
		ServerURL:              configuration.LDAP.ServerURL,
		BindDN:                 configuration.LDAP.BindDN,
		BindPw:                 configuration.LDAP.BindPw,
		MaxConnections:         configuration.LDAP.ConnectionPool.GetConnectionsLimit(),
		DialTimeout:            time.Duration(configuration.LDAP.ConnectionPool.GetDialTimeout()) * time.Second,
		RetryCount:             configuration.LDAP.ConnectionPool.GetRetryCount(),
		RetryDelay:             time.Duration(configuration.LDAP.ConnectionPool.GetRetryDelay()) * time.Second,
		ConnectionAliveTimeout: time.Duration(configuration.LDAP.ConnectionPool.GetConnectionAliveTimeout()) * time.Second,
	}

	ldapConnPool := backends.NewLdapConnectionPool(ldapConnPoolConfig)

	dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
		"ds_exporter",
		ldapConnPool,
		"cn=monitor",
		metrics.GetLdapServerMetrics(),
		prometheus.Labels{},
	),
	)

	dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
		"ds_exporter",
		ldapConnPool,
		"cn=snmp,cn=monitor",
		metrics.GetLdapServerSnmpMetrics(),
		prometheus.Labels{},
	),
	)

	for _, backend := range configuration.Global.Backends {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			ldapConnPool,
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
	if configuration.Global.BackendImplement == config.BackendBDB {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			ldapConnPool,
			"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
			metrics.GetLdapBDBServerCacheMetrics(),
			prometheus.Labels{},
		),
		)
	} else if configuration.Global.BackendImplement == config.BackendMDB {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			ldapConnPool,
			"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
			metrics.GetLdapMDBServerCacheMetrics(),
			prometheus.Labels{},
		),
		)
	}

	http.Handle(configuration.HTTP.GetMetricsPath(), promhttp.HandlerFor(dsMetricsRegistry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", defaultHttpResponse(configuration.HTTP.GetMetricsPath()))

	server := &http.Server{
		Addr:         configuration.HTTP.GetListenAddress(),
		Handler:      http.DefaultServeMux,
		ReadTimeout:  time.Duration(configuration.HTTP.GetReadTimeout()) * time.Second,
		WriteTimeout: time.Duration(configuration.HTTP.GetWriteTimeout()) * time.Second,
		IdleTimeout:  time.Duration(configuration.HTTP.GetIdleTimeout()) * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting HTTP server at %s", configuration.HTTP.GetListenAddress())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	signal := <-stop
	switch signal {
	case syscall.SIGINT:
		log.Println("SIGINT signal received")
	case syscall.SIGTERM:
		log.Println("SIGTREM signal received")
	}

	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown failed: %v", err)
	}

	if err = ldapConnPool.Close(ctx); err != nil {
		log.Printf("Error closing ldap pool: %v", err)
	}
	log.Println("Exporter stopped")
}
