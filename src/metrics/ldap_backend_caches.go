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
			Help:     "Some Help For This",
			Type:     prometheus.GaugeValue,
		},
		"dncachetries": {
			LdapName: "dncachetries",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
		"dncachehitratio": {
			LdapName: "dncachehitratio",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
		"currentdncachesize": {
			LdapName: "currentdncachesize",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
		"maxdncachesize": {
			LdapName: "maxdncachesize",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
		"currentdncachecount": {
			LdapName: "currentdncachecount",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
		"entrycachehits": {
			LdapName: "entrycachehits",
			Help:     "Some Help For This",
			Type:     prometheus.GaugeValue,
		},
		"entrycachetries": {
			LdapName: "entrycachetries",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
		"entrycachehitratio": {
			LdapName: "entrycachehitratio",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
		"currententrycachesize": {
			LdapName: "currententrycachesize",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
		"maxentrycachesize": {
			LdapName: "maxentrycachesize",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
		"currententrycachecount": {
			LdapName: "currententrycachecount",
			Help:     "Another Help For This",
			Type:     prometheus.GaugeValue,
		},
	}
}
