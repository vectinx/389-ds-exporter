package metrics

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/collectors"
	"389-ds-exporter/src/config"
	"389-ds-exporter/src/connections"
	"389-ds-exporter/src/utils"
)

func standardCollectors() []string {
	return []string{
		"server",
		"snmp-server",
		"ndn-cache",
		"ldbm-instance",
		"numsubordinates",
	}
}

func collectorEnabled(cfg *config.ExporterConfig, collector string) bool {
	if cfg.CollectorsDefault == "all" {
		return true
	}

	if cfg.CollectorsDefault == "none" &&
		slices.Contains(cfg.CollectorsEnabled, collector) {

		return true
	}

	if cfg.CollectorsDefault == "standard" {
		if slices.Contains(cfg.CollectorsEnabled, collector) ||
			slices.Contains(standardCollectors(), collector) {

			return true
		}
	}
	return false
}

// SetupPrometheusMetrics creates *prometheus.Registry, adds the required metrics and returns it.
func SetupPrometheusMetrics(
	cfg *config.ExporterConfig,
	connPool *connections.LdapConnectionPool,
	connPoolTimeout time.Duration,
) *prometheus.Registry {
	dsMetricsRegistry := prometheus.NewRegistry()

	dsCollector := collectors.NewDSCollector()

	if collectorEnabled(cfg, "server") {
		slog.Debug("Registering collector", "collector", "server")
		dsCollector.Register("server", collectors.NewLdapEntryCollector(
			"server",
			connPool,
			"cn=monitor",
			GetLdapServerMetrics(),
			prometheus.Labels{},
			connPoolTimeout,
		))
	}

	if collectorEnabled(cfg, "snmp-server") {
		slog.Debug("Registering collector", "collector", "snmp-server")
		dsCollector.Register("snmp-server", collectors.NewLdapEntryCollector(
			"snmp_server",
			connPool,
			"cn=snmp,cn=monitor",
			GetLdapServerSnmpMetrics(),
			prometheus.Labels{},
			connPoolTimeout,
		),
		)
	}

	if collectorEnabled(cfg, "numsubordinates") {
		slog.Debug("Registering collector", "collector", "numsubordinates")
		for _, entry := range cfg.NumSubordinateRecords {
			dsCollector.Register("numsubordinates_%s"+entry, collectors.NewLdapEntryCollector(
				"numsubordinates",
				connPool,
				entry,
				GetEntryCountAttr(),
				prometheus.Labels{"entry": entry},
				connPoolTimeout,
			),
			)
		}
	}

	if collectorEnabled(cfg, "ndn-cache") {
		slog.Debug("Registering collector", "collector", "ndn-cache")
		dsCollector.Register("ndn-cache", collectors.NewLdapEntryCollector(
			"ldbm",
			connPool,
			"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
			GetNdnCacheMetrics(),
			prometheus.Labels{},
			connPoolTimeout,
		),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connPoolTimeout)
	defer cancel()
	conn, err := connPool.Get(ctx)
	if err != nil {
		slog.Warn("Error obtaining connection to determine backend parameters", "err", err)
	} else {
		defer conn.Close()
	}

	backendInstances, err := utils.GetLdapBackendInstances(conn)
	if err != nil {
		slog.Warn("Error getting backend instances", "err", err)
	} else {
		if collectorEnabled(cfg, "ldbm-instance") {
			slog.Debug("Registering collector", "collector", "ldbm-instance")
			for _, instance := range backendInstances {
				slog.Info("Registeing metrics for backend instance", "instance", instance)
				dsCollector.Register("ldbm-instance_"+instance, collectors.NewLdapEntryCollector(
					"ldbm_instance",
					connPool,
					"cn=monitor,cn="+instance+",cn=ldbm database,cn=plugins,cn=config",
					GetLdapBackendCaches(),
					prometheus.Labels{"database": instance},
					connPoolTimeout,
				),
				)
			}
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
		switch *backendImplement {
		case config.BackendBDB:

			slog.Info("Berkeley DB backend implementation detected")

			if collectorEnabled(cfg, "bdb-caches") {
				slog.Debug("Registering collector", "collector", "bdb-caches")
				dsCollector.Register("bdb-caches", collectors.NewLdapEntryCollector(
					"bdb",
					connPool,
					"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
					GetLdapBDBServerCacheMetrics(),
					prometheus.Labels{},
					connPoolTimeout,
				),
				)
			}

			if collectorEnabled(cfg, "bdb-internal") {
				slog.Debug("Registering collector", "collector", "bdb-internal")
				dsCollector.Register("bdb-internal", collectors.NewLdapEntryCollector(
					"bdb",
					connPool,
					"cn=database,cn=monitor,cn=ldbm database,cn=plugins,cn=config",
					GetLdapBDBDatabaseLDBM(),
					prometheus.Labels{},
					connPoolTimeout,
				),
				)
			}
		case config.BackendMDB:
			slog.Info("LMDB backend implementation detected")

			if collectorEnabled(cfg, "lmdb-internal") {
				slog.Debug("Registering collector", "collector", "lmdb-internal")
				dsCollector.Register("lmdb-internal", collectors.NewLdapEntryCollector(
					"lmdb",
					connPool,
					"cn=database,cn=monitor,cn=ldbm database,cn=plugins,cn=config",
					GetLdapMDBDatabaseLDBM(),
					prometheus.Labels{},
					connPoolTimeout,
				),
				)
			}
		default:
			slog.Warn(
				"An unknown backend implementation type was detected. Backend metrics will not be collected",
				"backend",
				backendImplement,
			)
		}
	}

	dsMetricsRegistry.MustRegister(dsCollector)

	return dsMetricsRegistry
}
