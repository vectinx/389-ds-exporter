package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/collectors"
)

// GetLdapBDBServerCacheMetrics function returns map of
// specific for BDB backend attributes defining ldap server caches metrics.
func GetLdapBDBServerCacheMetrics() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"dbcachehits": {
			LdapName: "dbcachehits",
			Help:     "Number of requested pages found in the database",
			Type:     prometheus.GaugeValue,
		},
		"dbcachetries": {
			LdapName: "dbcachetries",
			Help:     "Total number of cache lookups",
			Type:     prometheus.GaugeValue,
		},
		"dbcachehitratio": {
			LdapName: "dbcachehitratio",
			Help:     "Percentage of requested pages found in the database cache",
			Type:     prometheus.GaugeValue,
		},
		"dbcachepagein": {
			LdapName: "dbcachepagein",
			Help:     "Number of pages read into the database cache",
			Type:     prometheus.GaugeValue,
		},
		"dbcachepageout": {
			LdapName: "dbcachepageout",
			Help:     "Number of pages written from the database cache to the backing file",
			Type:     prometheus.GaugeValue,
		},
		"dbcacheroevict": {
			LdapName: "dbcacheroevict",
			Help:     "Number of clean pages forced from the cache",
			Type:     prometheus.GaugeValue,
		},
		"dbcacherwevict": {
			LdapName: "dbcacherwevict",
			Help:     "Number of dirty pages forced from the cache",
			Type:     prometheus.GaugeValue,
		},
		"normalizeddncachetries": {
			LdapName: "normalizeddncachetries",
			Help:     "Total number of cache lookups since the instance was started",
			Type:     prometheus.GaugeValue,
		},
		"normalizeddncachehits": {
			LdapName: "normalizeddncachehits",
			Help:     "Normalized DNs found within the cache.",
			Type:     prometheus.GaugeValue,
		},
		"normalizeddncachemisses": {
			LdapName: "normalizeddncachemisses",
			Help:     "Normalized DNs not found within the cache",
			Type:     prometheus.GaugeValue,
		},
		"normalizeddncachehitratio": {
			LdapName: "normalizeddncachehitratio",
			Help:     "Percentage of the normalized DNs found in the cache",
			Type:     prometheus.GaugeValue,
		},
		"currentnormalizeddncachesize": {
			LdapName: "currentnormalizeddncachesize",
			Help:     "Current size of the normalized DN cache in bytes",
			Type:     prometheus.GaugeValue,
		},
		"maxnormalizeddncachesize": {
			LdapName: "maxnormalizeddncachesize",
			Help:     "Maximum size of NDn cache",
			Type:     prometheus.GaugeValue,
		},
		"currentnormalizeddncachecount": {
			LdapName: "currentnormalizeddncachecount",
			Help:     "Number of normalized cached DNs",
			Type:     prometheus.GaugeValue,
		},
	}
}

// GetLdapMDBServerCacheMetrics function returns a map
// of specific for MDB backend attributes defining ldap server caches metrics.
func GetLdapMDBServerCacheMetrics() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"normalizeddncachetries": {
			LdapName: "normalizeddncachetries",
			Help:     "Total number of cache lookups since the instance was started",
			Type:     prometheus.GaugeValue,
		},
		"normalizeddncachehits": {
			LdapName: "normalizeddncachehits",
			Help:     "Normalized DNs found within the cache.",
			Type:     prometheus.GaugeValue,
		},
		"normalizeddncachemisses": {
			LdapName: "normalizeddncachemisses",
			Help:     "Normalized DNs not found within the cache",
			Type:     prometheus.GaugeValue,
		},
		"normalizeddncachehitratio": {
			LdapName: "normalizeddncachehitratio",
			Help:     "Percentage of the normalized DNs found in the cache",
			Type:     prometheus.GaugeValue,
		},
		"currentnormalizeddncachesize": {
			LdapName: "currentnormalizeddncachesize",
			Help:     "Current size of the normalized DN cache in bytes",
			Type:     prometheus.GaugeValue,
		},
		"maxnormalizeddncachesize": {
			LdapName: "maxnormalizeddncachesize",
			Help:     "Maximum size of NDn cache",
			Type:     prometheus.GaugeValue,
		},
		"currentnormalizeddncachecount": {
			LdapName: "currentnormalizeddncachecount",
			Help:     "Number of normalized cached DNs",
			Type:     prometheus.GaugeValue,
		},
	}
}
