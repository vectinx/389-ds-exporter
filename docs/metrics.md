# Metrics collected by the exporter

To collect metrics, the exporter uses collectors. A *collector* is a component of the exporter responsible for collecting specific server metrics.
Below is a list of collectors available in the current version of `389-ds-exporter` and their purpose:
- [server](#server) - collects core server metrics.
- [snmp-server](#snmp-server) - collects extended server metrics. Complements but does not replace the `server` collector.
- [numsubordinates](#numsubordinates) - collects information about the number of entries in the configured DNs ([see config.md](config.md#global-settings)).
- [ndn-cache](#ndn-cache) - collects information about the usage of the Normalized DN Cache.
- [ldbm-instance](#ldbm-instance) - collects database backend metrics.
- [bdb-cache](#bdb-cache) - collects Berkeley DB cache information.
- [bdb-internal](#bdb-internal) - collects internal Berkeley DB metrics.
- [lmdb-internal](#lmdb-internal) - collects internal LMDB metrics.

Below is a detailed description of the metrics collected by each collector.

## `server`
The `server` collector collects core server metrics.</br>
Source: `cn=monitor`

### Metrics

#### ds_server_threads
Type: `gauge`
Attribute: `threads`
Current number of threads used to handle requests.

#### ds_server_currentconnections
Type: `gauge`
Attribute: `currentconnections`
Number of established connections.

---

#### ds_server_totalconnections
Type: `counter`
Attribute: `totalconnections`
Number of connections established by the server since startup.

---

#### ds_server_currentconnectionsatmaxthreads
Type: `gauge`
Attribute: `currentconnectionsatmaxthreads`
Current number of connections using the maximum allowed threads per connection.

---

#### ds_server_maxthreadsperconnhits
Type: `gauge`
Attribute: `maxthreadsperconnhits`
Shows how many times connections reached the thread-per-connection limit.

---

#### ds_server_dtablesize
Type: `gauge`
Attribute: `dtablesize`
Number of file descriptors available to the server.

---

#### ds_server_readwaiters
Type: `gauge`
Attribute: `readwaiters`
Number of connections with some requests waiting, not currently served by a server thread.

---

#### ds_server_opsinitiated
Type: `counter`
Attribute: `opsinitiated`
Number of operations initiated by the server since startup.

---

#### ds_server_opscompleted
Type: `counter`
Attribute: `opscompleted`
Number of operations completed by the server since startup.

---

#### ds_server_entriessent
Type: `counter`
Attribute: `entriessent`
Number of entries sent to clients since startup.

---

#### ds_server_bytessent
Type: `counter`
Attribute: `bytessent`
Number of bytes sent to clients since startup.

---

#### ds_server_nbackends
Type: `gauge`
Attribute: `nbackends`
Number of backends (databases, suffixes) served by the server.

---

#### ds_server_currenttime
Type: `gauge`
Attribute: `currenttime`
Current server time in UTC+0, Unix Timestamp format.

---

#### ds_server_starttime
Type: `gauge`
Attribute: `starttime`
Server startup time in UTC+0, Unix Timestamp format.

---

## `snmp-server`
The `snmp-server` collector collects extended server metrics. Complements but does not replace the `server` collector.</br>
Source: `cn=snmp,cn=monitor`

### Metrics

#### ds_snmp_server_anonymousbinds
Type: `counter`
Attribute: `anonymousbinds`
Number of anonymous BIND operations since server startup.

---

#### ds_snmp_server_unauthbinds
Type: `counter`
Attribute: `unauthbinds`
Number of unauthenticated BIND operations since server startup.

---

#### ds_snmp_server_simpleauthbinds
Type: `counter`
Attribute: `simpleauthbinds`
Number of simple BIND operations since server startup.

---

#### ds_snmp_server_strongauthbinds
Type: `counter`
Attribute: `strongauthbinds`
Number of strongauth BIND operations since server startup.

---

#### ds_snmp_server_bindsecurityerrors
Type: `counter`
Attribute: `bindsecurityerrors`
Number of times an incorrect password was specified in a BIND request.

---

#### ds_snmp_server_compareops
Type: `counter`
Attribute: `compareops`
Number of LDAP compare requests since server startup.

---

#### ds_snmp_server_addentryops
Type: `counter`
Attribute: `addentryops`
Number of LDAP add requests since server startup.

---

#### ds_snmp_server_removeentryops
Type: `counter`
Attribute: `removeentryops`
Number of LDAP delete requests since server startup.

---

#### ds_snmp_server_modifyentryops
Type: `counter`
Attribute: `modifyentryops`
Number of LDAP modify requests since server startup.

---

#### ds_snmp_server_modifyrdnops
Type: `counter`
Attribute: `modifyrdnops`
Number of LDAP modrdn requests since server startup.

---

#### ds_snmp_server_searchops
Type: `counter`
Attribute: `searchops`
Number of LDAP search requests since server startup.

---

#### ds_snmp_server_onelevelsearchops
Type: `counter`
Attribute: `onelevelsearchops`
Number of one-level search requests since server startup.

---

#### ds_snmp_server_wholesubtreesearchops
Type: `counter`
Attribute: `wholesubtreesearchops`
Number of subtree-level search requests since server startup.

---

#### ds_snmp_server_securityerrors
Type: `counter`
Attribute: `securityerrors`
Number of returned security-related errors, such as wrong passwords, invalid authentication methods, or higher security requirements.

---

#### ds_snmp_server_errors
Type: `counter`
Attribute: `errors`
Number of returned errors.

---

## `numsubordinates`
Collects information about the number of entries in the configured DNs.</br>
Source: configured DN

### Metrics

#### ds_numsubordinates_count
Type: `gauge`
Attribute: `numsubordinates`
Number of child entries in the DN.

---

## `ndn-cache`
Collects information about the usage of the Normalized DN Cache.</br>
Source: `cn=monitor,cn=ldbm database,cn=plugins,cn=config`

### Metrics

#### ds_ldbm_normalizeddncachetries
Type: `gauge`
Attribute: `normalizeddncachetries`
Total number of NDN cache accesses since server startup.

---

#### ds_ldbm_normalizeddncachehits
Type: `gauge`
Attribute: `normalizeddncachehits`
Number of normalized DNs found in cache since server startup.

---

#### ds_ldbm_normalizeddncachemisses
Type: `gauge`
Attribute: `normalizeddncachemisses`
Number of normalized DNs not found in cache since server startup.

---

#### ds_ldbm_normalizeddncachehitratio
Type: `gauge`
Attribute: `normalizeddncachehitratio`
Percentage of normalized DNs found in cache.

---

#### ds_ldbm_currentnormalizeddncachesize
Type: `gauge`
Attribute: `currentnormalizeddncachesize`
Current size of NDN cache in bytes.

---

#### ds_ldbm_maxnormalizeddncachesize
Type: `gauge`
Attribute: `maxnormalizeddncachesize`
Configured maximum size of NDN cache.

---

#### ds_ldbm_currentnormalizeddncachecount
Type: `gauge`
Attribute: `currentnormalizeddncachecount`
Number of normalized DNs currently cached.

---

## `ldbm-instance`
Collects database backend metrics. Database list is retrieved automatically at exporter startup.
Source: `cn=monitor,cn=<database_name>,cn=ldbm database,cn=plugins,cn=config`

### Metrics

#### ds_ldbm_instance_entrycachehits
Type: `counter`
Attribute: `entrycachehits`
Total number of successful entry cache accesses.

---

#### ds_ldbm_instance_entrycachetries
Type: `counter`
Attribute: `entrycachetries`
Total number of entry cache access attempts since server startup.

---

#### ds_ldbm_instance_entrycachehitratio
Type: `gauge`
Attribute: `entrycachehitratio`
Ratio of successful entry cache accesses to total attempts.

---

#### ds_ldbm_instance_currententrycachesize
Type: `gauge`
Attribute: `currententrycachesize`
Current size of entry cache in bytes.

---

#### ds_ldbm_instance_maxentrycachesize
Type: `gauge`
Attribute: `maxentrycachesize`
Maximum size of entry cache in bytes.

---

#### ds_ldbm_instance_currententrycachecount
Type: `gauge`
Attribute: `currententrycachecount`
Current number of entries stored in entry cache.

---

#### ds_ldbm_instance_dncachehits
Type: `counter`
Attribute: `dncachehits`
Number of accesses where the entry was found in DN cache.

---

#### ds_ldbm_instance_dncachetries
Type: `counter`
Attribute: `dncachetries`
Total number of DN cache accesses since server startup.

---

#### ds_ldbm_instance_dncachehitratio
Type: `gauge`
Attribute: `dncachehitratio`
Ratio of successful DN cache accesses to total accesses.

---

#### ds_ldbm_instance_currentdncachesize
Type: `gauge`
Attribute: `currentdncachesize`
Current size of DN cache in bytes.

---

#### ds_ldbm_instance_maxdncachesize
Type: `gauge`
Attribute: `maxdncachesize`
Maximum size of DN cache in bytes.

---

#### ds_ldbm_instance_currentdncachecount
Type: `gauge`
Attribute: `currentdncachecount`
Current number of entries in DN cache.

---

## `bdb-caches`
Collects Berkeley DB cache metrics.</br>
Source: `cn=monitor,cn=bdb,cn=ldbm database,cn=plugins,cn=config`

### Metrics

#### ds_bdb_cachehits
Type: `counter`
Attribute: `dbcachehits`
Number of pages found in the database cache without accessing disk files.

---

#### ds_bdb_cachetries
Type: `counter`
Attribute: `dbcachetries`
Total number of database cache accesses since server startup.

---

#### ds_bdb_cachehitratio
Type: `gauge`
Attribute: `dbcachehitratio`
Percentage of page requests found in the database cache.

---

#### ds_bdb_cachepagein
Type: `gauge`
Attribute: `dbcachepagein`
Number of pages loaded into the database cache from disk.

---

#### ds_bdb_cachepageout
Type: `gauge`
Attribute: `dbcachepageout`
Number of pages removed from the database cache to disk.

---

#### ds_bdb_cacheroevict
Type: `gauge`
Attribute: `dbcacheroevict`
Number of clean pages removed from cache.

---

#### ds_bdb_cacherwevict
Type: `gauge`
Attribute: `dbcacherwevict`
Number of dirty pages removed from cache.

---

## `bdb-internal`
Collects internal Berkeley DB metrics related to transactions, locks, pages, and transaction log.</br>
Source: `cn=monitor,cn=bdb,cn=ldbm database,cn=plugins,cn=config`

### Metrics

#### ds_bdb_abort_rate
Type: `counter`
Attribute: `nsslapd-db-abort-rate`
Number of aborted transactions.

---

#### ds_bdb_active_txns
Type: `gauge`
Attribute: `nsslapd-db-active-txns`
Number of transactions currently active and used by the database.

---

#### ds_bdb_cache_size_bytes
Type: `gauge`
Attribute: `nsslapd-db-cache-size-bytes`
Configured maximum size of database cache in bytes.

---

#### ds_bdb_cache_region_wait_rate
Type: `gauge`
Attribute: `nsslapd-db-cache-region-wait-rate`
Number of cases when a thread had to wait for a cache region lock.

---

#### ds_bdb_clean_pages
Type: `gauge`
Attribute: `nsslapd-db-clean-pages`
Number of clean pages in the database cache.

---

#### ds_bdb_commit_rate
Type: `counter`
Attribute: `nsslapd-db-commit-rate`
Number of committed transactions.

---

#### ds_bdb_deadlock_rate
Type: `gauge`
Attribute: `nsslapd-db-deadlock-rate`
Total number of detected deadlocks since server startup.

---

#### ds_bdb_dirty_pages
Type: `gauge`
Attribute: `nsslapd-db-dirty-pages`
Number of dirty pages in the database cache.

---

#### ds_bdb_hash_buckets
Type: `gauge`
Attribute: `nsslapd-db-hash-buckets`
Number of hash buckets in the buffer hash table.

---

#### ds_bdb_hash_elements_examine_rate
Type: `gauge`
Attribute: `nsslapd-db-hash-elements-examine-rate`
Number of hash elements examined during hash table searches.

---

#### ds_bdb_hash_search_rate
Type: `gauge`
Attribute: `nsslapd-db-hash-search-rate`
Number of searches in the buffer hash table.

---

#### ds_bdb_lock_conflicts
Type: `gauge`
Attribute: `nsslapd-db-lock-conflicts`
Number of times a lock could not be granted due to a conflict.

---

#### ds_bdb_lock_region_wait_rate
Type: `gauge`
Attribute: `nsslapd-db-lock-region-wait-rate`
Number of lock region waits.

---

#### ds_bdb_lock_request_rate
Type: `gauge`
Attribute: `nsslapd-db-lock-request-rate`
Total number of lock requests.

---

#### ds_bdb_lockers
Type: `gauge`
Attribute: `nsslapd-db-lockers`
Number of current lockers (subjects holding locks).

---

#### ds_bdb_configured_locks
Type: `gauge`
Attribute: `nsslapd-db-configured-locks`
Configured number of locks.

---

#### ds_bdb_current_locks
Type: `gauge`
Attribute: `nsslapd-db-current-locks`
Number of locks currently in use.

---

#### ds_bdb_max_locks
Type: `gauge`
Attribute: `nsslapd-db-max-locks`
Maximum number of locks used simultaneously since server startup.

---

#### ds_bdb_log_region_wait_rate
Type: `gauge`
Attribute: `nsslapd-db-log-region-wait-rate`
Number of transaction log region waits.

---

#### ds_bdb_log_write_rate
Type: `gauge`
Attribute: `nsslapd-db-log-write-rate`
Number of bytes written to the transaction log since the last checkpoint.

---

#### ds_bdb_longest_chain_length
Type: `gauge`
Attribute: `nsslapd-db-longest-chain-length`
Maximum chain length during hash table searches.

---

#### ds_bdb_page_create_rate
Type: `gauge`
Attribute: `nsslapd-db-page-create-rate`
Number of pages created in cache.

---

#### ds_bdb_page_read_rate
Type: `gauge`
Attribute: `nsslapd-db-page-read-rate`
Number of pages read into cache.

---

#### ds_bdb_page_ro_evict_rate
Type: `gauge`
Attribute: `nsslapd-db-page-ro-evict-rate`
Number of clean pages removed from cache.
> This duplicates `ds_bdb_cacheroevict`. Retained as in 389ds.

---

#### ds_bdb_page_rw_evict_rate
Type: `gauge`
Attribute: `nsslapd-db-page-rw-evict-rate`
Number of dirty pages removed from cache.
> This duplicates `ds_bdb_cacherwevict`. Retained as in 389ds.

---

#### ds_bdb_page_trickle_rate
Type: `gauge`
Attribute: `nsslapd-db-page-trickle-rate`
Number of dirty pages written using `memp_trickle`.

---

#### ds_bdb_page_write_rate
Type: `gauge`
Attribute: `nsslapd-db-page-write-rate`
Number of pages written from cache.

---

#### ds_bdb_pages_in_use
Type: `gauge`
Attribute: `nsslapd-db-pages-in-use`
Total number of pages (clean and dirty) currently used by cache.

---

#### ds_bdb_txn_region_wait_rate
Type: `gauge`
Attribute: `nsslapd-db-txn-region-wait-rate`
Number of transaction region wait cases.

---

#### ds_bdb_current_lock_objects
Type: `gauge`
Attribute: `nsslapd-db-current-lock-objects`
Current number of lock objects.

---

#### ds_bdb_max_lock_objects
Type: `gauge`
Attribute: `nsslapd-db-max-lock-objects`
Maximum number of lock objects recorded since server startup.

---

## `lmdb-internal`
Collects internal LMDB metrics related to transactions, files, and environment resources.</br>
Source: `cn=monitor,cn=mdb,cn=ldbm database,cn=plugins,cn=config`

### Metrics

#### ds_mdb_dbenvmapsize
Type: `gauge`
Attribute: `dbenvmapsize`
LMDB data file size in bytes.

---

#### ds_mdb_dbenvlastpageno
Type: `gauge`
Attribute: `dbenvlastpageno`
Number of pages used in the LMDB database file.

---

#### ds_mdb_dbenvlasttxnid
Type: `gauge`
Attribute: `dbenvlasttxnid`
Last LMDB transaction ID.

---

#### ds_mdb_dbenvmaxreaders
Type: `gauge`
Attribute: `dbenvmaxreaders`
Maximum number of allowed read threads in LMDB environment.

---

#### ds_mdb_dbenvnumreaders
Type: `gauge`
Attribute: `dbenvnumreaders`
Current number of read threads used in LMDB environment.

---

#### ds_mdb_dbenvnumdbis
Type: `gauge`
Attribute: `dbenvnumdbis`
Number of DBI (named databases) open in LMDB environment.

---

#### ds_mdb_waitingrwtxn
Type: `gauge`
Attribute: `waitingrwtxn`
Number of RW transactions currently waiting.

---

#### ds_mdb_activerwtxn
Type: `gauge`
Attribute: `activerwtxn`
Number of active RW transactions.

---

#### ds_mdb_abortrwtxn
Type: `gauge`
Attribute: `abortrwtxn`
Number of aborted RW transactions.

---

#### ds_mdb_commitrwtxn
Type: `gauge`
Attribute: `commitrwtxn`
Number of committed RW transactions.

---

#### ds_mdb_granttimerwtxn
Type: `gauge`
Attribute: `granttimerwtxn`
Description unavailable. If you know this attribute, please create an Issue in the project with a description or documentation link.

---

#### ds_mdb_lifetimerwtxn
Type: `gauge`
Attribute: `lifetimerwtxn`
Description unavailable. If you know this attribute, please create an Issue in the project with a description or documentation link.

---

#### ds_mdb_waitingrotxn
Type: `gauge`
Attribute: `waitingrotxn`
Number of RO transactions currently waiting.

---

#### ds_mdb_activerotxn
Type: `gauge`
Attribute: `activerotxn`
Number of active RO transactions.

---

#### ds_mdb_abortrotxn
Type: `gauge`
Attribute: `abortrotxn`
Number of aborted RO transactions.

---

#### ds_mdb_commitrotxn
Type: `gauge`
Attribute: `commitrotxn`
Number of committed RO transactions.

---

#### ds_mdb_granttimerotxn
Type: `gauge`
Attribute: `granttimerotxn`
Description unavailable. If you know this attribute, please create an Issue in the project with a description or documentation link.

---

#### ds_mdb_lifetimerotxn
Type: `gauge`
Attribute: `lifetimerotxn`
Description unavailable. If you know this attribute, please create an Issue in the project with a description or documentation link.
