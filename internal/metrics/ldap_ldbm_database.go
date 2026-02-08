/*
Package metrics provides ready-made sets of mappings of ldap attributes to prometheus metrics
*/
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"389-ds-exporter/internal/collectors"
)

// GetLdapBDBDatabaseLDBM function returns map of attributes defining specific ldap server ldbm database BDB metrics.
func GetLdapBDBDatabaseLDBM() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"txn_abort_total": {
			LdapName: "nsslapd-db-abort-rate",
			Help:     "Number of transactions that have been aborted.",
			Type:     prometheus.CounterValue,
		},
		"txn_active": {
			LdapName: "nsslapd-db-active-txns",
			Help:     "Number of transactions that are currently active (used by the database).",
			Type:     prometheus.GaugeValue,
		},
		"cache_size_bytes": {
			LdapName: "nsslapd-db-cache-size-bytes",
			Help:     "Total cache size in bytes.",
			Type:     prometheus.GaugeValue,
		},
		"cache_region_wait_total": {
			LdapName: "nsslapd-db-cache-region-wait-rate",
			Help:     "Number of times that a thread of control was forced to wait before obtaining the region lock.",
			Type:     prometheus.GaugeValue,
		},
		"cache_pages_clean": {
			LdapName: "nsslapd-db-clean-pages",
			Help:     "Clean pages currently in the cache.",
			Type:     prometheus.GaugeValue,
		},
		"txn_commit_total": {
			LdapName: "nsslapd-db-commit-rate",
			Help:     "Number of transactions that have been committed.",
			Type:     prometheus.CounterValue,
		},
		"deadlock_total": {
			LdapName: "nsslapd-db-deadlock-rate",
			Help:     "Number of deadlocks detected.",
			Type:     prometheus.GaugeValue,
		},
		"cache_pages_dirty": {
			LdapName: "nsslapd-db-dirty-pages",
			Help:     "Dirty pages currently in the cache.",
			Type:     prometheus.GaugeValue,
		},
		"cache_hash_buckets": {
			LdapName: "nsslapd-db-hash-buckets",
			Help:     "Number of hash buckets in buffer hash table.",
			Type:     prometheus.GaugeValue,
		},
		"cache_hash_elements_examined_total": {
			LdapName: "nsslapd-db-hash-elements-examine-rate",
			Help:     "Total number of hash elements traversed during hash table lookups.",
			Type:     prometheus.GaugeValue,
		},
		"cache_hash_lookups_total": {
			LdapName: "nsslapd-db-hash-search-rate",
			Help:     "Total number of buffer hash table lookups.",
			Type:     prometheus.GaugeValue,
		},
		"lock_conflicts_total": {
			LdapName: "nsslapd-db-lock-conflicts",
			Help:     "Total number of locks not immediately available due to conflicts.",
			Type:     prometheus.GaugeValue,
		},
		"lock_region_wait_total": {
			LdapName: "nsslapd-db-lock-region-wait-rate",
			Help:     "Number of times that a thread of control was forced to wait before obtaining the region lock.",
			Type:     prometheus.GaugeValue,
		},
		"lock_request_total": {
			LdapName: "nsslapd-db-lock-request-rate",
			Help:     "Total number of locks requested.",
			Type:     prometheus.GaugeValue,
		},
		"lockers": {
			LdapName: "nsslapd-db-lockers",
			Help:     "Number of current lockers.",
			Type:     prometheus.GaugeValue,
		},
		"locks_configured": {
			LdapName: "nsslapd-db-configured-locks",
			Help:     "Configured number of locks.",
			Type:     prometheus.GaugeValue,
		},
		"locks_current": {
			LdapName: "nsslapd-db-current-locks",
			Help:     "Number of locks currently used by the database.",
			Type:     prometheus.GaugeValue,
		},
		"locks_max": {
			LdapName: "nsslapd-db-max-locks",
			Help:     "The maximum number of locks at any one time.",
			Type:     prometheus.GaugeValue,
		},
		"log_region_wait_total": {
			LdapName: "nsslapd-db-log-region-wait-rate",
			Help:     "Number of times that a thread of control was forced to wait before obtaining the region lock.",
			Type:     prometheus.CounterValue,
		},
		"log_write_bytes_total": {
			LdapName: "nsslapd-db-log-write-rate",
			Help:     "Number of bytes written to the log since the last checkpoint.",
			Type:     prometheus.CounterValue,
		},
		"cache_hash_longest_chain": {
			LdapName: "nsslapd-db-longest-chain-length",
			Help:     "Longest chain ever encountered in buffer hash table lookups.",
			Type:     prometheus.GaugeValue,
		},
		"cache_page_create_total": {
			LdapName: "nsslapd-db-page-create-rate",
			Help:     "Pages created in the cache.",
			Type:     prometheus.CounterValue,
		},
		"cache_page_read_total": {
			LdapName: "nsslapd-db-page-read-rate",
			Help:     "Pages read into the cache.",
			Type:     prometheus.CounterValue,
		},
		"cache_page_ro_evict_total": {
			LdapName: "nsslapd-db-page-ro-evict-rate",
			Help:     "Clean pages forced from the cache.",
			Type:     prometheus.CounterValue,
		},
		"cache_page_rw_evict_total": {
			LdapName: "nsslapd-db-page-rw-evict-rate",
			Help:     "Dirty pages forced from the cache.",
			Type:     prometheus.CounterValue,
		},
		"cache_page_trickle_total": {
			LdapName: "nsslapd-db-page-trickle-rate",
			Help:     "Dirty pages written using the memp_trickle interface.",
			Type:     prometheus.CounterValue,
		},
		"cache_page_write_total": {
			LdapName: "nsslapd-db-page-write-rate",
			Help:     "Pages read into the cache.",
			Type:     prometheus.CounterValue,
		},
		"cache_pages_in_use": {
			LdapName: "nsslapd-db-pages-in-use",
			Help:     "All pages, clean or dirty, currently in use.",
			Type:     prometheus.GaugeValue,
		},
		"txn_region_wait_total": {
			LdapName: "nsslapd-db-txn-region-wait-rate",
			Help:     "Number of times that a thread of control was force to wait before obtaining the region lock.",
			Type:     prometheus.CounterValue,
		},
		"lock_objects_current": {
			LdapName: "nsslapd-db-current-lock-objects",
			Help:     "The number of current lock objects.",
			Type:     prometheus.GaugeValue,
		},
		"lock_objects_max": {
			LdapName: "nsslapd-db-max-lock-objects",
			Help:     "The maximum number of lock objects at any one time.",
			Type:     prometheus.GaugeValue,
		},
	}
}

// GetLdapMDBDatabaseLDBM function returns map of attributes defining specific ldap server ldbm database LMDB metrics.
func GetLdapMDBDatabaseLDBM() map[string]collectors.LdapMonitoredAttribute {
	return map[string]collectors.LdapMonitoredAttribute{
		"env_map_size_bytes": {
			LdapName: "dbenvmapsize",
			Help:     "LMDB Size of the data file in bytes.",
			Type:     prometheus.GaugeValue,
		},
		"env_last_page_number": {
			LdapName: "dbenvlastpageno",
			Help:     "LMDB Number of pages used",
			Type:     prometheus.GaugeValue,
		},
		"env_last_txn_id": {
			LdapName: "dbenvlasttxnid",
			Help:     "LMDB Last transaction ID",
			Type:     prometheus.GaugeValue,
		},
		"env_max_readers": {
			LdapName: "dbenvmaxreaders",
			Help:     "LMDB Max readers",
			Type:     prometheus.GaugeValue,
		},
		"env_num_readers": {
			LdapName: "dbenvnumreaders",
			Help:     "LMDB Number of readers used",
			Type:     prometheus.GaugeValue,
		},
		"env_num_dbis": {
			LdapName: "dbenvnumdbis",
			Help:     "Number of DBIs (named databases) within the LMDB environment",
			Type:     prometheus.GaugeValue,
		},
		"rw_txn_waiting": {
			LdapName: "waitingrwtxn",
			Help:     "LMDB Waiting RW transactions",
			Type:     prometheus.GaugeValue,
		},
		"rw_txn_active": {
			LdapName: "activerwtxn",
			Help:     "LMDB Active RW transactions",
			Type:     prometheus.GaugeValue,
		},
		"rw_txn_aborted": {
			LdapName: "abortrwtxn",
			Help:     "LMDB Aborted RW transactions",
			Type:     prometheus.GaugeValue,
		},
		"rw_txn_committed": {
			LdapName: "commitrwtxn",
			Help:     "LMDB Committed RW transactions",
			Type:     prometheus.GaugeValue,
		},
		"rw_txn_grant_time": {
			LdapName: "granttimerwtxn",
			Help:     "", // There is no clear explanation for this parameter yet.
			Type:     prometheus.GaugeValue,
		},
		"rw_txn_lifetime": {
			LdapName: "lifetimerwtxn",
			Help:     "", // There is no clear explanation for this parameter yet.
			Type:     prometheus.GaugeValue,
		},
		"ro_txn_waiting": {
			LdapName: "waitingrotxn",
			Help:     "MDB Waiting RO transactions",
			Type:     prometheus.GaugeValue,
		},
		"ro_txn_active": {
			LdapName: "activerotxn",
			Help:     "LMDB Active RO transactions",
			Type:     prometheus.GaugeValue,
		},
		"ro_txn_aborted": {
			LdapName: "abortrotxn",
			Help:     "LMDB Aborted RO transactions",
			Type:     prometheus.GaugeValue,
		},
		"ro_txn_committed": {
			LdapName: "commitrotxn",
			Help:     "LMDB Committed RO transactions",
			Type:     prometheus.GaugeValue,
		},
		"ro_txn_grant_time": {
			LdapName: "granttimerotxn",
			Help:     "", // There is no clear explanation for this parameter yet.
			Type:     prometheus.GaugeValue,
		},
		"ro_txn_lifetime": {
			LdapName: "lifetimerotxn",
			Help:     "", // There is no clear explanation for this parameter yet.
			Type:     prometheus.GaugeValue,
		},
	}
}
