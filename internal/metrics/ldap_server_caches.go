package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/internal/collectors"
)

// GetLdapBDBServerCacheMetrics function returns map of
// specific for BDB backend attributes defining ldap server caches metrics.
func GetLdapBDBServerCacheMetrics() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"dbcache_hits_total": {
			LdapName: "dbcachehits",
			Help:     "Number of requested pages found in the database",
			Type:     prometheus.CounterValue,
		},
		"dbcache_lookups_total": {
			LdapName: "dbcachetries",
			Help:     "Total number of cache lookups",
			Type:     prometheus.CounterValue,
		},
		"dbcache_hit_ratio": {
			LdapName: "dbcachehitratio",
			Help:     "Percentage of requested pages found in the database cache",
			Type:     prometheus.GaugeValue,
		},
		"dbcache_pages_in_total": {
			LdapName: "dbcachepagein",
			Help:     "Number of pages read into the database cache",
			Type:     prometheus.GaugeValue,
		},
		"dbcache_pages_out_total": {
			LdapName: "dbcachepageout",
			Help:     "Number of pages written from the database cache to the backing file",
			Type:     prometheus.GaugeValue,
		},
		"dbcache_evictions_clean_total": {
			LdapName: "dbcacheroevict",
			Help:     "Number of clean pages forced from the cache",
			Type:     prometheus.GaugeValue,
		},
		"dbcache_evictions_dirty_total": {
			LdapName: "dbcacherwevict",
			Help:     "Number of dirty pages forced from the cache",
			Type:     prometheus.GaugeValue,
		},
	}
}

// GetNdnCacheMetrics function returns a map
// of server NDN cache metrics.
func GetNdnCacheMetrics() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"ndn_cache_lookups_total": {
			LdapName: "normalizeddncachetries",
			Help:     "Total number of cache lookups since the server was started",
			Type:     prometheus.CounterValue,
		},
		"ndn_cache_hits_total": {
			LdapName: "normalizeddncachehits",
			Help:     "Normalized DNs found within the cache.",
			Type:     prometheus.CounterValue,
		},
		"ndn_cache_misses_total": {
			LdapName: "normalizeddncachemisses",
			Help:     "Normalized DNs not found within the cache",
			Type:     prometheus.CounterValue,
		},
		"ndn_cache_hit_ratio": {
			LdapName: "normalizeddncachehitratio",
			Help:     "Percentage of the normalized DNs found in the cache",
			Type:     prometheus.GaugeValue,
		},
		"ndn_cache_size_bytes": {
			LdapName: "currentnormalizeddncachesize",
			Help:     "Current size of the normalized DN cache in bytes",
			Type:     prometheus.GaugeValue,
		},
		"ndn_cache_max_size_bytes": {
			LdapName: "maxnormalizeddncachesize",
			Help:     "Maximum size of NDn cache",
			Type:     prometheus.GaugeValue,
		},
		"ndn_cache_entries": {
			LdapName: "currentnormalizeddncachecount",
			Help:     "Number of normalized cached DNs",
			Type:     prometheus.GaugeValue,
		},
	}
}
