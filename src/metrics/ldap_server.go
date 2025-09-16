package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/collectors"
)

// GetLdapServerMetrics is a map of attributes defining ldap server metrics.
func GetLdapServerMetrics() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"threads": {
			LdapName: "threads",
			Help:     "Current number of active threads used for handling requests",
			Type:     prometheus.CounterValue,
		},
		"currentconnections": {
			LdapName: "currentconnections",
			Help:     "Current established connections",
			Type:     prometheus.GaugeValue,
		},
		"totalconnections": {
			LdapName: "totalconnections",
			Help:     "Number of connections the server handles after it starts",
			Type:     prometheus.CounterValue,
		},
		"currentconnectionsatmaxthreads": {
			LdapName: "currentconnectionsatmaxthreads",
			Help:     "Number of connections currently utilizing the maximum allowed threads per connection",
			Type:     prometheus.GaugeValue,
		},
		"maxthreadsperconnhits": {
			LdapName: "maxthreadsperconnhits",
			Help:     "Displays how many times a connection hit max thread",
			Type:     prometheus.GaugeValue,
		},
		"dtablesize": {
			LdapName: "dtablesize",
			Help: `The number of file descriptors available to the directory.
Essentially, this value shows how many additional concurrent connections can be serviced by the directory`,
			Type: prometheus.GaugeValue,
		},
		"readwaiters": {
			LdapName: "readwaiters",
			Help:     "Number of threads waiting to read data from a client",
			Type:     prometheus.GaugeValue,
		},
		"opsinitiated": {
			LdapName: "opsinitiated",
			Help:     "Number of operations the server has initiated since it started",
			Type:     prometheus.GaugeValue,
		},
		"opscompleted": {
			LdapName: "opscompleted",
			Help:     "Number of operations the server has completed since it started.",
			Type:     prometheus.GaugeValue,
		},
		"entriessent": {
			LdapName: "entriessent",
			Help:     "Number of entries sent to clients since the server started",
			Type:     prometheus.GaugeValue,
		},
		"bytessent": {
			LdapName: "bytessent",
			Help:     "Number of bytes sent to clients after the server starts",
			Type:     prometheus.GaugeValue,
		},
		"nbackends": {
			LdapName: "nbackends",
			Help:     "Number of back ends (databases) the server services",
			Type:     prometheus.GaugeValue,
		},
		"currenttime": {
			LdapName: "currenttime",
			LdapType: collectors.Iso8601CompactString,
			Help:     "Current time of the server. The time is displayed in Greenwich Mean Time (GMT) in UTC format",
			Type:     prometheus.GaugeValue,
		},
		"starttime": {
			LdapName: "starttime",
			LdapType: collectors.Iso8601CompactString,
			Help:     "Time when the server started. The time is displayed in Greenwich Mean Time (GMT) in UTC format",
			Type:     prometheus.GaugeValue,
		},
	}
}
