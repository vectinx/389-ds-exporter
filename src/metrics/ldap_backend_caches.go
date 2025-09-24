/*
Package metrics provides ready-made sets of mappings of ldap attributes to prometheus metrics
*/
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/collectors"
)

// GetLdapBackendCaches function returns map of attributes defining specific ldap server backend metrics.
func GetLdapBackendCaches() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"dncachehits": {
			LdapName: "dncachehits",
			Help: `Number of times the server could process a request by obtaining a normalized
distinguished name (DN) from the DN cache rather than normalizing it again`,
			Type: prometheus.CounterValue,
		},
		"dncachetries": {
			LdapName: "dncachetries",
			Help:     "Total number of DN cache accesses since you started the instance",
			Type:     prometheus.CounterValue,
		},
		"dncachehitratio": {
			LdapName: "dncachehitratio",
			Help:     "Ratio of cache tries to successful DN cache hits. closer this value is to 100%, the better",
			Type:     prometheus.GaugeValue,
		},
		"currentdncachesize": {
			LdapName: "currentdncachesize",
			Help:     "Total size, in bytes, of DN currently present in the DN cache",
			Type:     prometheus.GaugeValue,
		},
		"maxdncachesize": {
			LdapName: "maxdncachesize",
			Help:     "Maximum size, in bytes, of DNs that DS can maintain in the DN cache",
			Type:     prometheus.GaugeValue,
		},
		"currentdncachecount": {
			LdapName: "currentdncachecount",
			Help:     "Number of DNs currently present in the DN cache",
			Type:     prometheus.GaugeValue,
		},
		"entrycachehits": {
			LdapName: "entrycachehits",
			Help:     "Total number of successful entry cache lookups",
			Type:     prometheus.CounterValue,
		},
		"entrycachetries": {
			LdapName: "entrycachetries",
			Help:     "Total number of entry cache lookups since you started the instance",
			Type:     prometheus.CounterValue,
		},
		"entrycachehitratio": {
			LdapName: "entrycachehitratio",
			Help:     "Number of entry cache tries to successful entry cache lookups",
			Type:     prometheus.GaugeValue,
		},
		"currententrycachesize": {
			LdapName: "currententrycachesize",
			Help:     "Total size, in bytes, of directory entries currently present in the entry cache",
			Type:     prometheus.GaugeValue,
		},
		"maxentrycachesize": {
			LdapName: "maxentrycachesize",
			Help:     "Maximum size, in bytes, of directory entries that {DS} can maintain in the entry cache",
			Type:     prometheus.GaugeValue,
		},
		"currententrycachecount": {
			LdapName: "currententrycachecount",
			Help:     "Current number of entries stored in the entry cache of a given backend",
			Type:     prometheus.GaugeValue,
		},
	}
}
