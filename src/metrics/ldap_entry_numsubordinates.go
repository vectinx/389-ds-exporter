/*
Package metrics provides ready-made sets of mappings of ldap attributes to prometheus metrics
*/
package metrics

import (
	"389-ds-exporter/src/collectors"

	"github.com/prometheus/client_golang/prometheus"
)

// GetEntryCountAttr function returns map of attributes defining specific ldap entry numsubordinates metric.
func GetEntryCountAttr() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"numsubordinates": {
			LdapName: "numsubordinates",
			Help:     "Indicates how many immediate subordinates an entry has.",
			Type:     prometheus.GaugeValue,
		},
	}
}
