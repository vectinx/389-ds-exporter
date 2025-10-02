package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/collectors"
	"389-ds-exporter/src/config"
	"389-ds-exporter/src/connections"
	"389-ds-exporter/src/utils"
)

// standardCollectors returns list of standard enabled collectors.
func standardCollectors() []string {
	return []string{
		"server",
		"snmp-server",
		"ndn-cache",
		"ldbm-instance",
		"numsubordinates",
	}
}

// registerCollectorIfEnabled registers the collector if it is enabled in the configuration.
func registerCollectorIfEnabled(
	dsCollector *collectors.DSCollector,
	name string,
	cfg *config.ExporterConfig,
	collector func() *collectors.LdapEntryCollector,
) {
	collectorEnabled := false
	switch cfg.CollectorsDefault {
	case "all":
		collectorEnabled = true
	case "none":
		collectorEnabled = slices.Contains(cfg.CollectorsEnabled, name)
	case "standard":
		collectorEnabled = slices.Contains(cfg.CollectorsEnabled, name) || slices.Contains(standardCollectors(), name)
	}

	if collectorEnabled {
		slog.Debug("Registering collector", "collector", name)
		dsCollector.Register(name, collector())
	}
}

// registerGeneralCollectors registers general collectors that are common to all backend types.
func registerGeneralCollectors(
	cfg *config.ExporterConfig,
	dsCollector *collectors.DSCollector,
	connPool *connections.LdapConnectionPool,
	connPoolTimeout time.Duration,
) {
	registerCollectorIfEnabled(dsCollector, "server", cfg, func() *collectors.LdapEntryCollector {
		return collectors.NewLdapEntryCollector(
			"server",
			connPool,
			"cn=monitor",
			GetLdapServerMetrics(),
			prometheus.Labels{},
			connPoolTimeout,
		)
	})

	registerCollectorIfEnabled(dsCollector, "snmp-server", cfg, func() *collectors.LdapEntryCollector {
		return collectors.NewLdapEntryCollector(
			"snmp_server",
			connPool,
			"cn=snmp,cn=monitor",
			GetLdapServerSnmpMetrics(),
			prometheus.Labels{},
			connPoolTimeout,
		)
	})

	registerCollectorIfEnabled(dsCollector, "ndn-cache", cfg, func() *collectors.LdapEntryCollector {
		return collectors.NewLdapEntryCollector(
			"ldbm",
			connPool,
			"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
			GetNdnCacheMetrics(),
			prometheus.Labels{},
			connPoolTimeout,
		)
	})

	for _, entry := range cfg.DSNumSubordinateRecords {
		e := entry
		registerCollectorIfEnabled(dsCollector, "numsubordinates_"+e, cfg, func() *collectors.LdapEntryCollector {
			return collectors.NewLdapEntryCollector(
				"numsubordinates",
				connPool,
				e,
				GetEntryCountAttr(),
				prometheus.Labels{"entry": entry},
				connPoolTimeout,
			)
		})
	}
}

// determineBackendInstances determines the list of backends
// to use based on the configuration and information in the LDAP directory.
func determineBackendInstances(cfg *config.ExporterConfig,
	pool *connections.LdapConnectionPool, timeout time.Duration) ([]string, error) {

	if len(cfg.DSBackendDBs) == 0 {
		slog.Debug("Backend instances not specified, detecting automatically")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		ldapConn, err := pool.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("error getting connection to determine backend parameters: %w", err)
		}

		defer ldapConn.Close()

		detectedInstances, err := utils.GetLdapBackendInstances(ldapConn)
		if err != nil {
			return nil, fmt.Errorf("backend instances detection error: %w", err)
		}

		slog.Debug("Using auto-detected backend instances")
		return detectedInstances, nil

	}

	slog.Debug("Using the backend instances specified in the configuration")
	return cfg.DSBackendDBs, nil
}

// determineBackendType determines backend type
// to use based on the configuration and information in the LDAP directory.
func determineBackendType(cfg *config.ExporterConfig,
	pool *connections.LdapConnectionPool, timeout time.Duration) (string, error) {

	if cfg.DSBackendType == "" {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		ldapConn, err := pool.Get(ctx)
		if err != nil {
			return "", fmt.Errorf("error getting connection to determine backend parameters: %w", err)
		}

		defer ldapConn.Close()

		slog.Debug("Backend type not specified, detecting automatically")

		detectedType, err := utils.GetLdapBackendType(ldapConn)
		if err != nil {
			return "", fmt.Errorf("backend type detection error: %w", err)
		}

		slog.Debug("Using auto-detected backend type")
		return *detectedType, nil

	}

	slog.Debug("Using the backend type specified in the configuration")
	return cfg.DSBackendType, nil
}

// SetupPrometheusMetrics creates *prometheus.Registry, adds the required metrics and returns it.
func SetupPrometheusMetrics(
	cfg *config.ExporterConfig,
	connPool *connections.LdapConnectionPool,
	connPoolTimeout time.Duration,
) *prometheus.Registry {

	slog.Info("Creating collectors...")
	defer slog.Info("Collectors created")
	dsMetricsRegistry := prometheus.NewRegistry()

	dsCollector := collectors.NewDSCollector()

	registerGeneralCollectors(cfg, dsCollector, connPool, connPoolTimeout)

	/*
		Since 389-ds has a different set of monitoring metrics for different backends (Berkley DB and LMDB),
		at the initialization stage we select the metrics that correspond to the selected backend
	*/

	backendType, err := determineBackendType(cfg, connPool, connPoolTimeout)

	if err != nil {
		slog.Error("Error detecting backend type", "err", err)
	} else {
		switch backendType {
		case config.BackendBDB:

			slog.Info("Berkeley DB backend implementation detected")

			registerCollectorIfEnabled(dsCollector, "bdb-caches", cfg, func() *collectors.LdapEntryCollector {
				return collectors.NewLdapEntryCollector(
					"bdb",
					connPool,
					"cn=monitor,cn=ldbm database,cn=plugins,cn=config",
					GetLdapBDBServerCacheMetrics(),
					prometheus.Labels{},
					connPoolTimeout,
				)
			})

			registerCollectorIfEnabled(dsCollector, "bdb-internal", cfg, func() *collectors.LdapEntryCollector {
				return collectors.NewLdapEntryCollector(
					"bdb",
					connPool,
					"cn=database,cn=monitor,cn=ldbm database,cn=plugins,cn=config",
					GetLdapBDBDatabaseLDBM(),
					prometheus.Labels{},
					connPoolTimeout,
				)
			})
		case config.BackendMDB:
			slog.Info("LMDB backend implementation detected")

			registerCollectorIfEnabled(dsCollector, "lmdb-internal", cfg, func() *collectors.LdapEntryCollector {
				return collectors.NewLdapEntryCollector(
					"lmdb",
					connPool,
					"cn=database,cn=monitor,cn=ldbm database,cn=plugins,cn=config",
					GetLdapMDBDatabaseLDBM(),
					prometheus.Labels{},
					connPoolTimeout,
				)
			})
		default:
			slog.Warn(
				"An unknown backend implementation type was detected. Backend metrics will not be collected",
				"backend",
				backendType,
			)
		}
	}

	detectedBackendInstances, err := determineBackendInstances(cfg, connPool, connPoolTimeout)
	if err != nil {
		slog.Error("Error detecting backend instances", "err", err)
	} else {
		for i := range detectedBackendInstances {
			registerCollectorIfEnabled(dsCollector, "ldbm-instance_"+detectedBackendInstances[i], cfg,
				func() *collectors.LdapEntryCollector {
					return collectors.NewLdapEntryCollector(
						"ldbm_instance",
						connPool,
						"cn=monitor,cn="+detectedBackendInstances[i]+",cn=ldbm database,cn=plugins,cn=config",
						GetLdapBackendCaches(),
						prometheus.Labels{"database": detectedBackendInstances[i]},
						connPoolTimeout,
					)
				})
		}
	}

	dsMetricsRegistry.MustRegister(dsCollector)

	return dsMetricsRegistry
}
