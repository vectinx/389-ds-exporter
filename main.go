package main

import (
	"fmt"
	"log"
	"net/http"
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
		listenAddress = kingpin.Flag("http.listen-address", "Address to listen on").
				Default(":9389").
				String()
		metricsPath = kingpin.Flag("http.metrics-path", "Path to expose metrics").
				Default("/metrics").
				String()
		httpReadTimeout = kingpin.Flag("http.read-timeout", "HTTP read timeout").
				Default("5").Int()
		httpWriteTimeout = kingpin.Flag("http.write-timeout", "HTTP write timeout").
					Default("10").Int()
		httpIdleimeout = kingpin.Flag("http.idle-timeout", "HTTP idle timeout").
				Default("120").Int()
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

	log.Printf("Configuration read successfuly")
	log.Printf("LDAP server URL: %v", configuration.LdapServerUrl)
	log.Printf("LDAP bind DN: %v", configuration.LdapBindDn)
	log.Printf("389-ds backend type: %v", configuration.BackendType)

	dsMetricsRegistry := prometheus.NewRegistry()

	ldapConnPoolConfig := backends.LdapConnectionPoolConfig{
		ServerURL:              configuration.LdapServerUrl,
		BindDN:                 configuration.LdapBindDn,
		BindPassword:           configuration.LdapBindPw,
		ConnectionsLimit:       1,
		MaxIdleTime:            600 * time.Second,
		MaxLifeTime:            3600 * time.Second,
		DialTimeout:            10 * time.Second,
		RetryCount:             3,
		RetryDelay:             2 * time.Second,
		ConnectionAliveTimeout: 2 * time.Second,
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

	for _, backend := range configuration.Backends {
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
		Since 389-ds has a different set of monitoring metrics for different backends (a and b),
		at the initialization stage we select the metrics that correspond to the selected backend
	*/
	if configuration.BackendType == config.BackendBDB {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			ldapConnPool,
			"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
			metrics.GetLdapBDBServerCacheMetrics(),
			prometheus.Labels{},
		),
		)
	} else if configuration.BackendType == config.BackendMDB {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			ldapConnPool,
			"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
			metrics.GetLdapMDBServerCacheMetrics(),
			prometheus.Labels{},
		),
		)
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(dsMetricsRegistry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", defaultHttpResponse(*metricsPath))

	server := &http.Server{
		Addr:         *listenAddress,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  time.Duration(*httpReadTimeout) * time.Second,
		WriteTimeout: time.Duration(*httpWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(*httpIdleimeout) * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
