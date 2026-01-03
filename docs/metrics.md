# Metrics collected by the exporter

To collect metrics, the exporter uses collectors. A *collector* is a component of the exporter responsible for collecting specific server metrics.
Below is a list of collectors available in the current version of 389-ds-exporter and their purposes:
- [server](#server) - collects basic server metrics.
- [snmp-server](#snmp-server) - collects extended server metrics. Complements but does not replace the `server` collector.
- [numsubordinates](#numsubordinates) - collects information about the number of entries in the DNs specified in the configuration ([see config.md](config.md#global-settings)).
- [ndn-cache](#ndn-cache) - collects information about the usage of the Normalized DN Cache.
- [ldbm-instance](#ldbm-instance) - collects backend database metrics.
- [bdb-cache](#bdb-cache) - collects information about Berkeley DB caches.
- [bdb-internal](#bdb-internal) - collects internal Berkeley DB metrics.
- [lmdb-internal](#lmdb-internal) - collects internal LMDB metrics.

Below is a detailed description of the metrics collected by each collector.

## `server`
The `server` collector collects basic server metrics</br>
Source: `cn=monitor`

#### ds_server_version

Type: `gauge`</br>
Attribute: `version`

A metric with a constant '1' value labeled by 389 Directory Server version.

#### ds_server_threads

Type: `gauge`</br>
Attribute: `threads`

The current number of threads used to process requests.

#### ds_server_connections

Type: `gauge`</br>
Attribute: `currentconnections`

The number of established connections.

#### ds_server_connections_total

Type: `counter`</br>
Attribute: `totalconnections`

The total number of connections established by the server since startup.

#### ds_server_connections_max_threads

Type: `gauge`</br>
Attribute: `currentconnectionsatmaxthreads`

The current number of connections using the maximum allowed number of threads per connection.

#### ds_server_max_threads_per_conn_hits_total

Type: `gauge`</br>
Attribute: `maxthreadsperconnhits`

Shows the number of times connections reached the thread limit per connection.

#### ds_server_dtablesize

Type: `gauge`</br>
Attribute: `dtablesize`

The number of file descriptors available to the server.

#### ds_server_read_waiters

Type: `gauge`</br>
Attribute: `readwaiters`

The number of connections with some requests in a waiting state and not currently being served by a thread on the server.

#### ds_server_ops_initiated_total

Type: `counter`</br>
Attribute: `opsinitiated`

The total number of operations initiated by the server since startup.

#### ds_server_ops_completed_total

Type: `counter`</br>
Attribute: `opscompleted`

The total number of operations completed by the server since startup.

#### ds_server_entries_sent_total

Type: `counter`</br>
Attribute: `entriessent`

The total number of entries sent to clients since startup.

#### ds_server_bytes_sent_total

Type: `counter`</br>
Attribute: `bytessent`

The total number of bytes sent to clients since startup.

#### ds_server_backends

Type: `gauge`</br>
Attribute: `nbackends`

The number of backends (databases, suffixes) served by the server.

#### ds_server_current_time_seconds

Type: `gauge`</br>
Attribute: `currenttime`

The current server time in UTC+0 timezone in Unix Timestamp format.

#### ds_server_start_time_seconds

Type: `gauge`</br>
Attribute: `starttime`

The server startup time in UTC+0 timezone in Unix Timestamp format.

## `snmp-server`
The `snmp-server` collector collects extended server metrics. It complements but does not replace the `server` collector</br>
Metrics source: `cn=snmp,cn=monitor`

#### ds_snmp_server_bind_anonymous_total

Type: `counter`</br>
Attribute: `anonymousbinds`

The total number of anonymous BIND operations since server startup.

#### ds_snmp_server_bind_unauth_total

Type: `counter`</br>
Attribute: `unauthbinds`

The total number of unauthenticated BIND operations since server startup.

#### ds_snmp_server_bind_simple_total

Type: `counter`</br>
Attribute: `simpleauthbinds`

The total number of simple authentication BIND operations since server startup.

#### ds_snmp_server_bind_strong_total

Type: `counter`</br>
Attribute: `strongauthbinds`

The total number of strong authentication BIND operations since server startup.

#### ds_snmp_server_bind_security_errors_total

Type: `counter`</br>
Attribute: `bindsecurityerrors`

The number of times an incorrect password was provided in a BIND request.

#### ds_snmp_server_compare_operations_total

Type: `counter`</br>
Attribute: `compareops`

The total number of LDAP compare requests since server startup.

#### ds_snmp_server_add_operations_total

Type: `counter`</br>
Attribute: `addentryops`

The total number of LDAP add requests since server startup.

#### ds_snmp_server_delete_operations_total

Type: `counter`</br>
Attribute: `removeentryops`

The total number of LDAP delete requests since server startup.

#### ds_snmp_server_modify_operations_total

Type: `counter`</br>
Attribute: `modifyentryops`

The total number of LDAP modify requests since server startup.

#### ds_snmp_server_modify_rdn_operations_total

Type: `counter`</br>
Attribute: `modifyrdnops`

The total number of LDAP modrdn requests since server startup.

#### ds_snmp_server_search_operations_total

Type: `counter`</br>
Attribute: `searchops`

The total number of LDAP search requests since server startup.

#### ds_snmp_server_search_onelevel_operations_total

Type: `counter`</br>
Attribute: `onelevelsearchops`

The total number of one-level search requests since server startup.

#### ds_snmp_server_search_whole_subtree_operations_total

Type: `counter`</br>
Attribute: `wholesubtreesearchops`

The total number of subtree-level search requests since server startup.

#### ds_snmp_server_security_errors_total

Type: `counter`</br>
Attribute: `securityerrors`

The number of security-related errors returned, such as incorrect passwords, incorrect authentication methods, or requirements for higher security levels.

#### ds_snmp_server_errors_total

Type: `counter`</br>
Attribute: `errors`

The total number of errors returned.

## `numsubordinates`
The `numsubordinates` collector collects information about the number of entries in the DNs specified in the configuration.</br>
Metrics source: DNs specified in the configuration

#### ds_numsubordinates_count

Type: `gauge`</br>
Attribute: `numsubordinates`

The number of child DN entries.

## `ndn-cache`
The `ndn-cache` collector collects information about the usage of the Normalized DN Cache.</br>
Source: `cn=monitor,cn=ldbm database,cn=plugins,cn=config`

#### ds_ldbm_ndn_cache_lookups_total

Type: `counter`</br>
Attribute: `normalizeddncachetries`

The total number of lookups to the NDN cache since server startup.

#### ds_ldbm_ndn_cache_hits_total

Type: `counter`</br>
Attribute: `normalizeddncachehits`

The number of normalized DNs found in the cache since server startup.

#### ds_ldbm_ndn_cache_misses_total

Type: `counter`</br>
Attribute: `normalizeddncachemisses`

The number of normalized DNs not found in the cache since server startup.

#### ds_ldbm_ndn_cache_hit_ratio

Type: `gauge`</br>
Attribute: `normalizeddncachehitratio`

The percentage of normalized DNs found in the cache.

#### ds_ldbm_ndn_cache_size_bytes

Type: `gauge`</br>
Attribute: `currentnormalizeddncachesize`

The current size of the NDN cache in bytes.

#### ds_ldbm_ndn_cache_max_size_bytes

Type: `gauge`</br>
Attribute: `maxnormalizeddncachesize`

The configured maximum size of the NDN cache.

#### ds_ldbm_ndn_cache_entries

Type: `gauge`</br>
Attribute: `currentnormalizeddncachecount`

The number of normalized DNs currently cached.

## `ldbm-instance`
The `ldbm-instance` collector collects backend database metrics. The list of databases is obtained automatically when the exporter starts.</br>
Metrics source: `cn=monitor,cn=<database name>,cn=ldbm database,cn=plugins,cn=config`

#### ds_ldbm_instance_entry_cache_hits_total

Type: `counter`</br>
Attribute: `entrycachehits`

The total number of successful entry cache hits.

#### ds_ldbm_instance_entry_cache_lookups_total

Type: `counter`</br>
Attribute: `entrycachetries`

The total number of attempts to access the entry cache since server startup.

#### ds_ldbm_instance_entry_cache_hit_ratio

Type: `gauge`</br>
Attribute: `entrycachehitratio`

The ratio of successful entry cache hits to the total number of attempts.

#### ds_ldbm_instance_entry_cache_size_bytes

Type: `gauge`</br>
Attribute: `currententrycachesize`

The current size of the entry cache in bytes.

#### ds_ldbm_instance_entry_cache_max_size_bytes

Type: `gauge`</br>
Attribute: `maxentrycachesize`

The maximum size of the entry cache in bytes.

#### ds_ldbm_instance_entry_cache_count

Type: `gauge`</br>
Attribute: `currententrycachecount`

The current number of entries stored in the entry cache.

#### ds_ldbm_instance_dn_cache_hits_total

Type: `counter`</br>
Attribute: `dncachehits`

The number of hits when an entry was found in the cache.

#### ds_ldbm_instance_dn_cache_lookups_total

Type: `counter`</br>
Attribute: `dncachetries`

The total number of accesses to the DN cache since server startup.

#### ds_ldbm_instance_dn_cache_hit_ratio

Type: `gauge`</br>
Attribute: `dncachehitratio`

The ratio of successful DN cache hits to the total number of accesses.

#### ds_ldbm_instance_dn_cache_size_bytes

Type: `gauge`</br>
Attribute: `currentdncachesize`

The current size of the DN cache in bytes.

#### ds_ldbm_instance_dn_cache_max_size_bytes

Type: `gauge`</br>
Attribute: `maxdncachesize`

The maximum size of the DN cache in bytes.

#### ds_ldbm_instance_dn_cache_count

Type: `gauge`</br>
Attribute: `currentdncachecount`

The current number of entries in the DN cache.

## `bdb-cache`
The `bdb-cache` collector collects BerkeleyDB cache metrics.</br>
Metrics source: `cn=monitor,cn=bdb,cn=ldbm database,cn=plugins,cn=config`

#### ds_bdb_dbcache_hits_total

Type: `counter`</br>
Attribute: `dbcachehits`

The number of pages found in the database cache without accessing disk files.

#### ds_bdb_dbcache_lookups_total

Type: `counter`</br>
Attribute: `dbcachetries`

The total number of cache accesses since server startup.

#### ds_bdb_dbcache_hit_ratio

Type: `gauge`</br>
Attribute: `dbcachehitratio`

The percentage of page requests found in the database cache. Higher values indicate more efficient cache usage.

#### ds_bdb_dbcache_pages_in_total

Type: `gauge`</br>
Attribute: `dbcachepagein`

The number of pages loaded into the database cache from disk.

#### ds_bdb_dbcache_pages_out_total

Type: `gauge`</br>
Attribute: `dbcachepageout`

The number of pages evicted from the database cache to disk.

#### ds_bdb_dbcache_evictions_clean_total

Type: `gauge`</br>
Attribute: `dbcacheroevict`

The number of "clean" (no disk write required) pages evicted from the cache.

#### ds_bdb_dbcache_evictions_dirty_total

Type: `gauge`</br>
Attribute: `dbcacherwevict`

The number of "dirty" (disk write required) pages evicted from the cache.

## `bdb-internal`
The `bdb-internal` collector collects internal BerkeleyDB metrics related to transactions, locks, pages, and transaction logs.</br>
Metrics source: `cn=monitor,cn=bdb,cn=ldbm database,cn=plugins,cn=config`

#### ds_bdb_txn_abort_total

Type: `counter`</br>
Attribute: `nsslapd-db-abort-rate`

The number of aborted transactions.

#### ds_bdb_txn_active

Type: `gauge`</br>
Attribute: `nsslapd-db-active-txns`

The number of transactions currently active and in use by the database.

#### ds_bdb_cache_size_bytes

Type: `gauge`</br>
Attribute: `nsslapd-db-cache-size-bytes`

The maximum configured size of the database cache in bytes.

#### ds_bdb_cache_region_wait_total

Type: `gauge`</br>
Attribute: `nsslapd-db-cache-region-wait-rate`

The number of times a thread had to wait to acquire a cache region lock.

#### ds_bdb_cache_pages_clean

Type: `gauge`</br>
Attribute: `nsslapd-db-clean-pages`

The number of "clean" pages in the database cache.

#### ds_bdb_txn_commit_total

Type: `counter`</br>
Attribute: `nsslapd-db-commit-rate`

The number of committed transactions.

#### ds_bdb_deadlock_total

Type: `gauge`</br>
Attribute: `nsslapd-db-deadlock-rate`

The total number of deadlocks detected since server startup.

#### ds_bdb_cache_pages_dirty

Type: `gauge`</br>
Attribute: `nsslapd-db-dirty-pages`

The number of "dirty" pages in the database cache.

#### ds_bdb_cache_hash_buckets

Type: `gauge`</br>
Attribute: `nsslapd-db-hash-buckets`

The number of hash buckets in the buffer hash table.

#### ds_bdb_cache_hash_elements_examined_total

Type: `gauge`</br>
Attribute: `nsslapd-db-hash-elements-examine-rate`

The number of hash elements examined during searches in the hash table.

#### ds_bdb_cache_hash_lookups_total

Type: `gauge`</br>
Attribute: `nsslapd-db-hash-search-rate`

The number of searches in the buffer hash table.

#### ds_bdb_lock_conflicts_total

Type: `gauge`</br>
Attribute: `nsslapd-db-lock-conflicts`

The number of times a lock could not be granted due to a conflict.

#### ds_bdb_lock_region_wait_total

Type: `gauge`</br>
Attribute: `nsslapd-db-lock-region-wait-rate`

The number of region lock wait cases.

#### ds_bdb_lock_request_total

Type: `gauge`</br>
Attribute: `nsslapd-db-lock-request-rate`

The total number of lock request attempts.

#### ds_bdb_lockers

Type: `gauge`</br>
Attribute: `nsslapd-db-lockers`

The number of current "lockers" (entities holding locks).

#### ds_bdb_locks_configured

Type: `gauge`</br>
Attribute: `nsslapd-db-configured-locks`

The configured number of locks.

#### ds_bdb_locks_current

Type: `gauge`</br>
Attribute: `nsslapd-db-current-locks`

The number of locks currently in use.

#### ds_bdb_locks_max

Type: `gauge`</br>
Attribute: `nsslapd-db-max-locks`

The maximum number of locks used simultaneously since server startup.

#### ds_bdb_log_region_wait_total

Type: `gauge`</br>
Attribute: `nsslapd-db-log-region-wait-rate`

The number of waits for transaction log region locks.

#### ds_bdb_log_write_bytes_total

Type: `gauge`</br>
Attribute: `nsslapd-db-log-write-rate`

The number of bytes written to the log since the last transaction log checkpoint.

#### ds_bdb_cache_hash_longest_chain

Type: `gauge`</br>
Attribute: `nsslapd-db-longest-chain-length`

The maximum length of a chain during a search in the buffer hash table.

#### ds_bdb_cache_page_create_total

Type: `counter`</br>
Attribute: `nsslapd-db-page-create-rate`

The number of pages created in the cache.

#### ds_bdb_cache_page_read_total

Type: `counter`</br>
Attribute: `nsslapd-db-page-read-rate`

The number of pages read into the cache.

#### ds_bdb_cache_page_ro_evict_total

Type: `counter`</br>
Attribute: `nsslapd-db-page-ro-evict-rate`

The number of "clean" pages evicted from the cache.
> This value duplicates the `ds_bdb_cacheroevict` metric. It would be possible to keep only one of these metrics,
but 389ds duplicates it for some reason, so let's keep it here as well.

#### ds_bdb_cache_page_rw_evict_total

Type: `counter`</br>
Attribute: `nsslapd-db-page-rw-evict-rate`

The number of "dirty" pages evicted from the cache.
> This value duplicates the `ds_bdb_cacherwevict` metric. It would be possible to keep only one of these metrics,
but 389ds duplicates it for some reason, so let's keep it here as well.

#### ds_bdb_cache_page_trickle_total

Type: `counter`</br>
Attribute: `nsslapd-db-page-trickle-rate`

The number of "dirty" pages written using the `memp_trickle` interface.

#### ds_bdb_cache_page_write_total

Type: `counter`</br>
Attribute: `nsslapd-db-page-write-rate`

The number of pages written from the cache.

#### ds_bdb_cache_pages_in_use

Type: `gauge`</br>
Attribute: `nsslapd-db-pages-in-use`

The total number of pages (both clean and dirty) currently in use by the cache.

#### ds_bdb_txn_region_wait_total

Type: `counter`</br>
Attribute: `nsslapd-db-txn-region-wait-rate`

The number of waits for transaction region locks.

#### ds_bdb_lock_objects_current

Type: `gauge`</br>
Attribute: `nsslapd-db-current-lock-objects`

The current number of lock objects.

#### ds_bdb_lock_objects_max

Type: `gauge`</br>
Attribute: `nsslapd-db-max-lock-objects`

The maximum number of lock objects recorded since server startup.

## `lmdb-internal`
The `lmdb-internal` collector collects internal LMDB metrics related to transactions, files, and environment resources.</br>
Metrics source: `cn=monitor,cn=mdb,cn=ldbm database,cn=plugins,cn=config`

#### ds_lmdb_env_map_size_bytes

Type: `gauge`</br>
Attribute: `dbenvmapsize`

The size of the LMDB data file in bytes.

#### ds_lmdb_env_last_page_number

Type: `gauge`</br>
Attribute: `dbenvlastpageno`

The number of pages used in the LMDB database file.

#### ds_lmdb_env_last_txn_id

Type: `gauge`</br>
Attribute: `dbenvlasttxnid`

The ID of the last LMDB transaction.

#### ds_lmdb_env_max_readers

Type: `gauge`</br>
Attribute: `dbenvmaxreaders`

The maximum number of reader threads allowed in the LMDB environment.

#### ds_lmdb_env_num_readers

Type: `gauge`</br>
Attribute: `dbenvnumreaders`

The current number of reader threads in use in the LMDB environment.

#### ds_lmdb_env_num_dbis

Type: `gauge`</br>
Attribute: `dbenvnumdbis`

The number of DBI (named databases) opened in the LMDB environment.

#### ds_lmdb_rw_txn_waiting

Type: `gauge`</br>
Attribute: `waitingrwtxn`

The number of RW (read/write) transactions waiting.

#### ds_lmdb_rw_txn_active

Type: `gauge`</br>
Attribute: `activerwtxn`

The number of active RW (read/write) transactions.

#### ds_lmdb_rw_txn_aborted

Type: `gauge`</br>
Attribute: `abortrwtxn`

The number of aborted RW transactions.

#### ds_lmdb_rw_txn_committed

Type: `gauge`</br>
Attribute: `commitrwtxn`

The number of successfully committed RW transactions.

#### ds_lmdb_rw_txn_grant_time

Type: `gauge`</br>
Attribute: `granttimerwtxn`

Description unavailable. Documentation could not be found for this attribute. If you know what this attribute means, please create an Issue in the project with a description or documentation link.

#### ds_lmdb_rw_txn_lifetime

Type: `gauge`</br>
Attribute: `lifetimerwtxn`

Description unavailable. Documentation could not be found for this attribute. If you know what this attribute means, please create an Issue in the project with a description or documentation link.

#### ds_lmdb_ro_txn_waiting

Type: `gauge`</br>
Attribute: `waitingrotxn`

The number of RO (read-only) transactions waiting.

#### ds_lmdb_ro_txn_active

Type: `gauge`</br>
Attribute: `activerotxn`

The number of active RO transactions.

#### ds_lmdb_ro_txn_aborted

Type: `gauge`</br>
Attribute: `abortrotxn`

The number of aborted RO transactions.

#### ds_lmdb_ro_txn_committed

Type: `gauge`</br>
Attribute: `commitrotxn`

The number of successfully committed RO transactions.

#### ds_lmdb_ro_txn_grant_time

Type: `gauge`</br>
Attribute: `granttimerotxn`

Description unavailable. Documentation could not be found for this attribute. If you know what this attribute means, please create an Issue in the project with a description or documentation link.

#### ds_lmdb_ro_txn_lifetime

Type: `gauge`</br>
Attribute: `lifetimerotxn`

Description unavailable. Documentation could not be found for this attribute. If you know what this attribute means, please create an Issue in the project with a description or documentation link.
