package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/internal/collectors"
)

// GetLdapServerSnmpMetrics returns a map of attributes defining ldap server backend metrics.
func GetLdapServerSnmpMetrics() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"bind_anonymous_total": {
			LdapName: "anonymousbinds",
			Help:     "Number of anonymous bind requests",
			Type:     prometheus.CounterValue,
		},
		"bind_unauth_total": {
			LdapName: "unauthbinds",
			Help:     "Number of unauthenticated (anonymous) binds",
			Type:     prometheus.CounterValue,
		},
		"bind_simple_total": {
			LdapName: "simpleauthbinds",
			Help:     "Number of LDAP simple bind requests (DN and password)",
			Type:     prometheus.CounterValue,
		},
		"bind_strong_total": {
			LdapName: "strongauthbinds",
			Help:     "Number of LDAP SASL bind requests, for all SASL mechanisms",
			Type:     prometheus.CounterValue,
		},
		"bind_security_errors_total": {
			LdapName: "bindsecurityerrors",
			Help:     "Number of times an invalid password was given in a bind request.",
			Type:     prometheus.CounterValue,
		},
		"compare_operations_total": {
			LdapName: "compareops",
			Help:     "Number of LDAP compare requests",
			Type:     prometheus.CounterValue,
		},
		"add_operations_total": {
			LdapName: "addentryops",
			Help:     "Number of LDAP add requests.",
			Type:     prometheus.CounterValue,
		},
		"delete_operations_total": {
			LdapName: "removeentryops",
			Help:     "Number of LDAP delete requests",
			Type:     prometheus.CounterValue,
		},
		"modify_operations_total": {
			LdapName: "modifyentryops",
			Help:     "Number of LDAP modify requests",
			Type:     prometheus.CounterValue,
		},
		"modify_rdn_operations_total": {
			LdapName: "modifyrdnops",
			Help:     "Number of LDAP modify RDN (modrdn) requests",
			Type:     prometheus.CounterValue,
		},
		"search_operations_total": {
			LdapName: "searchops",
			Help:     "Number of LDAP search requests",
			Type:     prometheus.CounterValue,
		},
		"search_onelevel_operations_total": {
			LdapName: "onelevelsearchops",
			Help:     "Number of one-level search operations",
			Type:     prometheus.CounterValue,
		},
		"search_whole_subtree_operations_total": {
			LdapName: "wholesubtreesearchops",
			Help:     "Number of subtree-level search operations",
			Type:     prometheus.CounterValue,
		},
		"security_errors_total": {
			LdapName: "securityerrors",
			Help: `Number of errors returned that were security related, such as invalid passwords,
unknown or invalid authentication methods, or stronger authentication required`,
			Type: prometheus.CounterValue,
		},
		"errors_total": {
			LdapName: "errors",
			Help:     "Number of errors returned",
			Type:     prometheus.CounterValue,
		},
	}
}
