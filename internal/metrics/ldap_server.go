package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/internal/collectors"
)

// GetLdapServerMetrics is a map of attributes defining ldap server metrics.
func GetLdapServerMetrics() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"version": {
			LdapName: "version",
			Help:     "A metric with a constant '1' value labeled by 389 Directory Server version",
			Type:     prometheus.GaugeValue,
			LdapType: collectors.StringLabel,
		},
		"threads": {
			LdapName: "threads",
			Help:     "Current number of active threads used for handling requests",
			Type:     prometheus.GaugeValue,
		},
		"connections": {
			LdapName: "currentconnections",
			Help:     "Current established connections",
			Type:     prometheus.GaugeValue,
		},
		"connections_total": {
			LdapName: "totalconnections",
			Help:     "Number of connections the server handles after it starts",
			Type:     prometheus.CounterValue,
		},
		"connections_max_threads": {
			LdapName: "currentconnectionsatmaxthreads",
			Help:     "Number of connections currently utilizing the maximum allowed threads per connection",
			Type:     prometheus.GaugeValue,
		},
		"max_threads_per_conn_hits_total": {
			LdapName: "maxthreadsperconnhits",
			Help:     "Displays how many times a connection hit max thread",
			Type:     prometheus.GaugeValue,
		},
		"dtablesize": {
			LdapName: "dtablesize",
			Help:     "The number of file descriptors available to the directory",
			Type:     prometheus.GaugeValue,
		},
		"read_waiters": {
			LdapName: "readwaiters",
			Help: `Number of connections where some requests are pending
and not currently being serviced by a thread in Directory Server`,
			Type: prometheus.GaugeValue,
		},
		"ops_initiated_total": {
			LdapName: "opsinitiated",
			Help:     "Number of operations the server has initiated since it started",
			Type:     prometheus.CounterValue,
		},
		"ops_completed_total": {
			LdapName: "opscompleted",
			Help:     "Number of operations the server has completed since it started.",
			Type:     prometheus.CounterValue,
		},
		"entries_sent_total": {
			LdapName: "entriessent",
			Help:     "Number of entries sent to clients since the server started",
			Type:     prometheus.CounterValue,
		},
		"bytes_sent_total": {
			LdapName: "bytessent",
			Help:     "Number of bytes sent to clients after the server starts",
			Type:     prometheus.CounterValue,
		},
		"backends": {
			LdapName: "nbackends",
			Help:     "Number of back ends (databases) the server services",
			Type:     prometheus.GaugeValue,
		},
		"current_time_seconds": {
			LdapName: "currenttime",
			LdapType: collectors.Iso8601CompactString,
			Help:     "Current time of the server. The time is displayed in Greenwich Mean Time (GMT) in UTC format",
			Type:     prometheus.CounterValue,
		},
		"start_time_seconds": {
			LdapName: "starttime",
			LdapType: collectors.Iso8601CompactString,
			Help:     "Time when the server started. The time is displayed in Greenwich Mean Time (GMT) in UTC format",
			Type:     prometheus.GaugeValue,
		},
	}
}
