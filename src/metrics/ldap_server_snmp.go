package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/collectors"
)

// LdapServerSnmpMetrics returns a map of attributes defining ldap server backend metrics.
func GetLdapServerSnmpMetrics() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"anonymousbinds": {
			LdapName: "anonymousbinds",
			Help:     "Number of anonymous bind requests",
			Type:     prometheus.CounterValue,
		},
		"unauthbinds": {
			LdapName: "unauthbinds",
			Help:     "Number of unauthenticated (anonymous) binds",
			Type:     prometheus.CounterValue,
		},
		"simpleauthbinds": {
			LdapName: "simpleauthbinds",
			Help:     "Number of LDAP simple bind requests (DN and password)",
			Type:     prometheus.CounterValue,
		},
		"strongauthbinds": {
			LdapName: "strongauthbinds",
			Help:     "Number of LDAP SASL bind requests, for all SASL mechanisms",
			Type:     prometheus.CounterValue,
		},
		"bindsecurityerrors": {
			LdapName: "bindsecurityerrors",
			Help:     "Number of number of times an invalid password was given in a bind request.",
			Type:     prometheus.CounterValue,
		},
		"compareops": {
			LdapName: "compareops",
			Help:     "Number of LDAP compare requests",
			Type:     prometheus.CounterValue,
		},
		"addentryops": {
			LdapName: "addentryops",
			Help:     "Number of LDAP add requests.",
			Type:     prometheus.CounterValue,
		},
		"removeentryops": {
			LdapName: "removeentryops",
			Help:     "Number of LDAP delete requests",
			Type:     prometheus.CounterValue,
		},
		"modifyentryops": {
			LdapName: "modifyentryops",
			Help:     "Number of LDAP modify requests",
			Type:     prometheus.CounterValue,
		},
		"modifyrdnops": {
			LdapName: "modifyrdnops",
			Help:     "Number of LDAP modify RDN (modrdn) requests",
			Type:     prometheus.CounterValue,
		},
		"searchops": {
			LdapName: "searchops",
			Help:     "Number of LDAP search requests",
			Type:     prometheus.CounterValue,
		},
		"onelevelsearchops": {
			LdapName: "onelevelsearchops",
			Help:     "Number of one-level search operations",
			Type:     prometheus.CounterValue,
		},
		"wholesubtreesearchops": {
			LdapName: "wholesubtreesearchops",
			Help:     "Number of subtree-level search operations",
			Type:     prometheus.CounterValue,
		},
		"securityerrors": {
			LdapName: "securityerrors",
			Help: `Number of errors returned that were security related, such as invalid passwords,
unknown or invalid authentication methods, or stronger authentication required`,
			Type: prometheus.CounterValue,
		},
		"errors": {
			LdapName: "errors",
			Help:     "Number of errors returned",
			Type:     prometheus.CounterValue,
		},
	}
}
