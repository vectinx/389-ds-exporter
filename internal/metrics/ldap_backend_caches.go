/*
Package metrics provides ready-made sets of mappings of ldap attributes to prometheus metrics
*/
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/internal/collectors"
)

// GetLdapBackendCaches function returns map of attributes defining specific ldap server backend metrics.
func GetLdapBackendCaches() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"dn_cache_hits_total": {
			LdapName: "dncachehits",
			Help: `Number of times the server could process a request by obtaining a normalized
distinguished name (DN) from the DN cache rather than normalizing it again`,
			Type: prometheus.CounterValue,
		},
		"dn_cache_lookups_total": {
			LdapName: "dncachetries",
			Help:     "Total number of DN cache accesses since you started the instance",
			Type:     prometheus.CounterValue,
		},
		"dn_cache_hit_ratio": {
			LdapName: "dncachehitratio",
			Help:     "Ratio of cache tries to successful DN cache hits. closer this value is to 100%, the better",
			Type:     prometheus.GaugeValue,
		},
		"dn_cache_size_bytes": {
			LdapName: "currentdncachesize",
			Help:     "Total size, in bytes, of DN currently present in the DN cache",
			Type:     prometheus.GaugeValue,
		},
		"dn_cache_max_size_bytes": {
			LdapName: "maxdncachesize",
			Help:     "Maximum size, in bytes, of DNs that DS can maintain in the DN cache",
			Type:     prometheus.GaugeValue,
		},
		"dn_cache_count": {
			LdapName: "currentdncachecount",
			Help:     "Number of DNs currently present in the DN cache",
			Type:     prometheus.GaugeValue,
		},
		"entry_cache_hits_total": {
			LdapName: "entrycachehits",
			Help:     "Total number of successful entry cache lookups",
			Type:     prometheus.CounterValue,
		},
		"entry_cache_lookups_total": {
			LdapName: "entrycachetries",
			Help:     "Total number of entry cache lookups since you started the instance",
			Type:     prometheus.CounterValue,
		},
		"entry_cache_hit_ratio": {
			LdapName: "entrycachehitratio",
			Help:     "Number of entry cache tries to successful entry cache lookups",
			Type:     prometheus.GaugeValue,
		},
		"entry_cache_size_bytes": {
			LdapName: "currententrycachesize",
			Help:     "Total size, in bytes, of directory entries currently present in the entry cache",
			Type:     prometheus.GaugeValue,
		},
		"entry_cache_max_size_bytes": {
			LdapName: "maxentrycachesize",
			Help:     "Maximum size, in bytes, of directory entries that {DS} can maintain in the entry cache",
			Type:     prometheus.GaugeValue,
		},
		"entry_cache_count": {
			LdapName: "currententrycachecount",
			Help:     "Current number of entries stored in the entry cache of a given backend",
			Type:     prometheus.GaugeValue,
		},
	}
}
