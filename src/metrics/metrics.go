package metrics

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/collectors"
	"389-ds-exporter/src/config"
	"389-ds-exporter/src/connections"
	"389-ds-exporter/src/utils"
)

// SetupPrometheusMetrics creates *prometheus.Registry, adds the required metrics and returns it.
func SetupPrometheusMetrics(
	cfg *config.ExporterConfiguration,
	connPool *connections.LdapConnectionPool,
	connPoolTimeout time.Duration,
) *prometheus.Registry {
	dsMetricsRegistry := prometheus.NewRegistry()

	dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
		"ds_exporter",
		connPool,
		"cn=monitor",
		GetLdapServerMetrics(),
		prometheus.Labels{},
		connPoolTimeout,
	),
	)

	dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
		"ds_exporter",
		connPool,
		"cn=snmp,cn=monitor",
		GetLdapServerSnmpMetrics(),
		prometheus.Labels{},
		connPoolTimeout,
	),
	)

	for _, entry := range cfg.Global.NumSubordinatesRecords {
		dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
			"ds_exporter",
			connPool,
			entry,
			GetEntryCountAttr(),
			prometheus.Labels{"entry": entry},
			connPoolTimeout,
		),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := connPool.Get(ctx)
	if err != nil {
		slog.Warn("Error obtaining connection to determine backend parameters", "err", err)
		return dsMetricsRegistry
	}

	defer conn.Close()

	backendInstances, err := utils.GetLdapBackendInstances(conn)
	if err != nil {
		slog.Warn("Error getting backend instances", "err", err)
	} else {
		for _, instance := range backendInstances {
			slog.Info("Registeing metrics for backend instance", "instance", instance)
			dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
				"ds_exporter",
				connPool,
				"cn=monitor,cn="+instance+",cn=ldbm database,cn=plugins,cn=config",
				GetLdapBackendCaches(),
				prometheus.Labels{"database": instance},
				connPoolTimeout,
			),
			)
		}
	}

	/*
		Since 389-ds has a different set of monitoring metrics for different backends (Berkley DB and LMDB),
		at the initialization stage we select the metrics that correspond to the selected backend
	*/
	backendImplement, err := utils.GetLdapBackendType(conn)
	if err != nil {
		slog.Warn("Error getting backend implementation type", "err", err)
	} else {
		switch backendImplement {
		case config.BackendBDB:
			slog.Info("Berkeley DB backend implementation detected")
			dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
				"ds_exporter",
				connPool,
				"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
				GetLdapBDBServerCacheMetrics(),
				prometheus.Labels{},
				connPoolTimeout,
			),
			)
			dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
				"ds_exporter",
				connPool,
				"cn=database,cn=monitor,cn=ldbm database,cn=plugins,cn=config",
				GetLdapBDBDatabaseLDBM(),
				prometheus.Labels{},
				connPoolTimeout,
			),
			)
		case config.BackendMDB:
			slog.Info("LMDB backend implementation detected")
			dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
				"ds_exporter",
				connPool,
				"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
				GetLdapMDBServerCacheMetrics(),
				prometheus.Labels{},
				connPoolTimeout,
			),
			)
			dsMetricsRegistry.MustRegister(collectors.NewLdapEntryCollector(
				"ds_exporter",
				connPool,
				"cn=database,cn=monitor,cn=ldbm database,cn=plugins,cn=config",
				GetLdapMDBDatabaseLDBM(),
				prometheus.Labels{},
				connPoolTimeout,
			),
			)
		default:
			slog.Warn(
				"An unknown backend implementation type was detected. Backend metrics will not be collected",
				"backend",
				backendImplement,
			)
		}
	}

	return dsMetricsRegistry
}
