/*
Package metrics provides ready-made sets of mappings of ldap attributes to prometheus metrics
*/
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/src/collectors"
)

// GetLdapBDBDatabaseLDBM function returns map of attributes defining specific ldap server ldbm database BDB metrics.
func GetLdapBDBDatabaseLDBM() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"db_abort_rate": {
			LdapName: "nsslapd-db-abort-rate",
			Help:     "Number of transactions that have been aborted.",
			Type:     prometheus.GaugeValue,
		},
		"db_active_txns": {
			LdapName: "nsslapd-db-active-txns",
			Help:     "Number of transactions that are currently active (used by the database).",
			Type:     prometheus.GaugeValue,
		},
		"db_cache_size_bytes": {
			LdapName: "nsslapd-db-cache-size-bytes",
			Help:     "Total cache size in bytes.",
			Type:     prometheus.GaugeValue,
		},
		"db_cache_region_wait_rate": {
			LdapName: "nsslapd-db-cache-region-wait-rate",
			Help:     "Number of times that a thread of control was forced to wait before obtaining the region lock.",
			Type:     prometheus.GaugeValue,
		},
		"db_clean_pages": {
			LdapName: "nsslapd-db-clean-pages",
			Help:     "Clean pages currently in the cache.",
			Type:     prometheus.GaugeValue,
		},
		"db_commits_rate": {
			LdapName: "nsslapd-db-commit-rate",
			Help:     "Number of transactions that have been committed.",
			Type:     prometheus.GaugeValue,
		},
		"db_deadlock_rate": {
			LdapName: "nsslapd-db-deadlock-rate",
			Help:     "Number of deadlocks detected.",
			Type:     prometheus.GaugeValue,
		},
		"db_dirty_pages": {
			LdapName: "nsslapd-db-dirty-pages",
			Help:     "Dirty pages currently in the cache.",
			Type:     prometheus.GaugeValue,
		},
		"db_hash_buckets": {
			LdapName: "nsslapd-db-hash-buckets",
			Help:     "Number of hash buckets in buffer hash table.",
			Type:     prometheus.GaugeValue,
		},
		"db_hash_elements_examine_rate": {
			LdapName: "nsslapd-db-hash-elements-examine-rate",
			Help:     "Total number of hash elements traversed during hash table lookups.",
			Type:     prometheus.GaugeValue,
		},
		"db_hash_search_rate": {
			LdapName: "nsslapd-db-hash-search-rate",
			Help:     "Total number of buffer hash table lookups.",
			Type:     prometheus.GaugeValue,
		},
		"db_lock_conflicts": {
			LdapName: "nsslapd-db-lock-conflicts",
			Help:     "Total number of locks not immediately available due to conflicts.",
			Type:     prometheus.GaugeValue,
		},
		"db_lock_region_wait_rate": {
			LdapName: "nsslapd-db-lock-region-wait-rate",
			Help:     "Number of times that a thread of control was forced to wait before obtaining the region lock.",
			Type:     prometheus.GaugeValue,
		},
		"db_lock_request_rate": {
			LdapName: "nsslapd-db-lock-request-rate",
			Help:     "Total number of locks requested.",
			Type:     prometheus.GaugeValue,
		},
		"db_lockers": {
			LdapName: "nsslapd-db-lockers",
			Help:     "Number of current lockers.",
			Type:     prometheus.GaugeValue,
		},
		"db_configured_locks": {
			LdapName: "nsslapd-db-configured-locks",
			Help:     "Configured number of locks.",
			Type:     prometheus.GaugeValue,
		},
		"db_current_locks": {
			LdapName: "nsslapd-db-current-locks",
			Help:     "Number of locks currently used by the database.",
			Type:     prometheus.GaugeValue,
		},
		"db_max_locks": {
			LdapName: "nsslapd-db-max-locks",
			Help:     "The maximum number of locks at any one time.",
			Type:     prometheus.GaugeValue,
		},
		"db_log_region_wait_rate": {
			LdapName: "nsslapd-db-log-region-wait-rate",
			Help:     "Number of times that a thread of control was forced to wait before obtaining the region lock.",
			Type:     prometheus.GaugeValue,
		},
		"db_log_write_rate": {
			LdapName: "nsslapd-db-log-write-rate",
			Help:     "Number of bytes written to the log since the last checkpoint.",
			Type:     prometheus.GaugeValue,
		},
		"db_longest_chain_length": {
			LdapName: "nsslapd-db-longest-chain-length",
			Help:     "Longest chain ever encountered in buffer hash table lookups.",
			Type:     prometheus.GaugeValue,
		},
		"db_page_create_rate": {
			LdapName: "nsslapd-db-page-create-rate",
			Help:     "Pages created in the cache.",
			Type:     prometheus.GaugeValue,
		},
		"db_page_read_rate": {
			LdapName: "nsslapd-db-page-read-rate",
			Help:     "Pages read into the cache.",
			Type:     prometheus.GaugeValue,
		},
		"db_page_ro_evict_rate": {
			LdapName: "nsslapd-db-page-ro-evict-rate",
			Help:     "Clean pages forced from the cache.",
			Type:     prometheus.GaugeValue,
		},
		"db_page_rw_evict_rate": {
			LdapName: "nsslapd-db-page-rw-evict-rate",
			Help:     "Dirty pages forced from the cache.",
			Type:     prometheus.GaugeValue,
		},
		"db_page_trickle_rate": {
			LdapName: "nsslapd-db-page-trickle-rate",
			Help:     "Dirty pages written using the memp_trickle interface.",
			Type:     prometheus.GaugeValue,
		},
		"db_page_write_rate": {
			LdapName: "nsslapd-db-page-write-rate",
			Help:     "Pages read into the cache.",
			Type:     prometheus.GaugeValue,
		},
		"db_pages_in_use": {
			LdapName: "nsslapd-db-pages-in-use",
			Help:     "All pages, clean or dirty, currently in use.",
			Type:     prometheus.GaugeValue,
		},
		"db_txn_region_wait_rate": {
			LdapName: "nsslapd-db-txn-region-wait-rate",
			Help:     "Number of times that a thread of control was force to wait before obtaining the region lock.",
			Type:     prometheus.GaugeValue,
		},
		"db_current_lock_objects": {
			LdapName: "nsslapd-db-current-lock-objects",
			Help:     "The number of current lock objects.",
			Type:     prometheus.GaugeValue,
		},
		"db_max_lock_objects": {
			LdapName: "nsslapd-db-max-lock-objects",
			Help:     "The maximum number of lock objects at any one time.",
			Type:     prometheus.GaugeValue,
		},
	}
}

// GetLdapMDBDatabaseLDBM function returns map of attributes defining specific ldap server ldbm database LMDB metrics.
func GetLdapMDBDatabaseLDBM() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"dbenvmapsize": {
			LdapName: "dbenvmapsize",
			Help:     "LMDB Size of the data file in bytes.",
			Type:     prometheus.GaugeValue,
		},
		"dbenvlastpageno": {
			LdapName: "dbenvlastpageno",
			Help:     "LMDB Number of pages used",
			Type:     prometheus.GaugeValue,
		},
		"dbenvlasttxnid": {
			LdapName: "dbenvlasttxnid",
			Help:     "LMDB Last transaction ID",
			Type:     prometheus.GaugeValue,
		},
		"dbenvmaxreaders": {
			LdapName: "dbenvmaxreaders",
			Help:     "LMDB Max readers",
			Type:     prometheus.GaugeValue,
		},
		"dbenvnumreaders": {
			LdapName: "dbenvnumreaders",
			Help:     "LMDB Number of readers used",
			Type:     prometheus.GaugeValue,
		},
		"dbenvnumdbis": {
			LdapName: "dbenvnumdbis",
			Help:     "",
			Type:     prometheus.GaugeValue,
		},
		"waitingrwtxn": {
			LdapName: "waitingrwtxn",
			Help:     "LMDB Waiting RW transactions",
			Type:     prometheus.GaugeValue,
		},
		"activerwtxn": {
			LdapName: "activerwtxn",
			Help:     "LMDB Active RW transactions",
			Type:     prometheus.GaugeValue,
		},
		"abortrwtxn": {
			LdapName: "abortrwtxn",
			Help:     "LMDB Aborted RW transactions",
			Type:     prometheus.GaugeValue,
		},
		"commitrwtxn": {
			LdapName: "commitrwtxn",
			Help:     "LMDB Commited RW transactions",
			Type:     prometheus.GaugeValue,
		},
		"granttimerwtxn": {
			LdapName: "granttimerwtxn",
			Help:     "",
			Type:     prometheus.GaugeValue,
		},
		"lifetimerwtxn": {
			LdapName: "lifetimerwtxn",
			Help:     "",
			Type:     prometheus.GaugeValue,
		},
		"waitingrotxn": {
			LdapName: "waitingrotxn",
			Help:     "MDB Waiting RO transactions",
			Type:     prometheus.GaugeValue,
		},
		"activerotxn": {
			LdapName: "activerotxn",
			Help:     "LMDB Active RO transactions",
			Type:     prometheus.GaugeValue,
		},
		"abortrotxn": {
			LdapName: "abortrotxn",
			Help:     "LMDB Aborted RO transactions",
			Type:     prometheus.GaugeValue,
		},
		"commitrotxn": {
			LdapName: "commitrotxn",
			Help:     "LMDB Commited RO transactions",
			Type:     prometheus.GaugeValue,
		},
		"granttimerotxn": {
			LdapName: "granttimerotxn",
			Help:     "",
			Type:     prometheus.GaugeValue,
		},
		"lifetimerotxn": {
			LdapName: "lifetimerotxn",
			Help:     "",
			Type:     prometheus.GaugeValue,
		},
	}
}
