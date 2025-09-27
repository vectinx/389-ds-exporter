/*
Package metrics provides ready-made sets of mappings of ldap attributes to prometheus metrics
*/
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/collectors"
)

// GetEntryCountAttr function returns map of attributes defining specific ldap entry numsubordinates metric.
func GetEntryCountAttr() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"count": {
			LdapName: "numsubordinates",
			Help:     "Indicates how many immediate subordinates an entry has.",
			Type:     prometheus.GaugeValue,
		},
	}
}
