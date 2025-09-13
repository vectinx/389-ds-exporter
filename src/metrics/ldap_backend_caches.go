/*
Package metrics provides ready-made sets of mappings of ldap attributes to prometheus metrics
*/
package metrics

import (
	"389-ds-exporter/src/collectors"

	"github.com/prometheus/client_golang/prometheus"
)

// GetLdapBackendCaches function returns map of attributes defining specific ldap server backend metrics.
func GetLdapBackendCaches() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"dncachehits": {
			LdapName: "dncachehits",
			Help: `The number of times the server could process a request by
obtaining a normalized distinguished name (DN) from the DN cache rather than normalizing it again`,
			Type: prometheus.GaugeValue,
		},
		"dncachetries": {
			LdapName: "dncachetries",
			Help:     "The total number of DN cache accesses since you started the instance",
			Type:     prometheus.GaugeValue,
		},
		"dncachehitratio": {
			LdapName: "dncachehitratio",
			Help:     "The ratio of cache tries to successful DN cache hits. The closer this value is to 100%, the better",
			Type:     prometheus.GaugeValue,
		},
		"currentdncachesize": {
			LdapName: "currentdncachesize",
			Help:     "The total size, in bytes, of DN currently present in the DN cache",
			Type:     prometheus.GaugeValue,
		},
		"maxdncachesize": {
			LdapName: "maxdncachesize",
			Help:     "The maximum size, in bytes, of DNs that DS can maintain in the DN cache",
			Type:     prometheus.GaugeValue,
		},
		"currentdncachecount": {
			LdapName: "currentdncachecount",
			Help:     "The number of DNs currently present in the DN cache",
			Type:     prometheus.GaugeValue,
		},
		"entrycachehits": {
			LdapName: "entrycachehits",
			Help:     "The total number of successful entry cache lookups",
			Type:     prometheus.GaugeValue,
		},
		"entrycachetries": {
			LdapName: "entrycachetries",
			Help:     "The total number of entry cache lookups since you started the instance",
			Type:     prometheus.GaugeValue,
		},
		"entrycachehitratio": {
			LdapName: "entrycachehitratio",
			Help:     "The number of entry cache tries to successful entry cache lookups",
			Type:     prometheus.GaugeValue,
		},
		"currententrycachesize": {
			LdapName: "currententrycachesize",
			Help:     "The total size, in bytes, of directory entries currently present in the entry cache",
			Type:     prometheus.GaugeValue,
		},
		"maxentrycachesize": {
			LdapName: "maxentrycachesize",
			Help:     "The maximum size, in bytes, of directory entries that {DS} can maintain in the entry cache",
			Type:     prometheus.GaugeValue,
		},
		"currententrycachecount": {
			LdapName: "currententrycachecount",
			Help:     "The current number of entries stored in the entry cache of a given backend",
			Type:     prometheus.GaugeValue,
		},
	}
}
