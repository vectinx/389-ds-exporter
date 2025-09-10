package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

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
				Default(":9389").
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

	configuration := config.ReadConfig(*configFilePath)
	log.Println(configuration)

	dsMetricsRegistry := prometheus.NewRegistry()

	dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
		"ds_exporter",
		configuration.LdapServerUrl,
		configuration.LdapBindDn,
		configuration.LdapBindPw,
		"cn=monitor",
		metrics.GetLdapServerMetrics(),
		prometheus.Labels{},
	),
	)

	dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
		"ds_exporter",
		configuration.LdapServerUrl,
		configuration.LdapBindDn,
		configuration.LdapBindPw,
		"cn=snmp,cn=monitor",
		metrics.GetLdapServerSnmpMetrics(),
		prometheus.Labels{},
	),
	)

	dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
		"ds_exporter",
		configuration.LdapServerUrl,
		configuration.LdapBindDn,
		configuration.LdapBindPw,
		"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
		metrics.GetLdapServerCacheMetrics(),
		prometheus.Labels{},
	),
	)

	for _, backend := range configuration.Backends {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			configuration.LdapServerUrl,
			configuration.LdapBindDn,
			configuration.LdapBindPw,
			"cn=monitor,cn="+backend+",cn=ldbm database,cn=plugins,cn=config",
			metrics.GetLdapBackendCaches(),
			prometheus.Labels{"database": backend},
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
