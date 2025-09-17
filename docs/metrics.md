# Metrics collected by exporter

## cn=monitor

Metrics collected from `cn=monitor`

| Metric                                      | Attribute                       | Description |
| ------------------------------------------- | ------------------------------- | ----------- |
| ds_exporter_threads                         | threads                         | Current number of active threads used for handling requests  |
| ds_exporter_currentconnections              | currentconnections              | Current established connections  |
| ds_exporter_totalconnections                | totalconnections                | Number of connections the server handles after it starts  |
| ds_exporter_currentconnectionsatmaxthreads  | currentconnectionsatmaxthreads  | Number of connections currently utilizing the maximum allowed threads per connection  |
| ds_exporter_maxthreadsperconnhits           | maxthreadsperconnhits           | Displays how many times a connection hit max thread  |
| ds_exporter_dtablesize                      | dtablesize                      | The number of file descriptors available to the directory  |
| ds_exporter_readwaiters                     | readwaiters                     | Number of connections where some requests are pending and not currently being serviced by a thread in Directory Server  |
| ds_exporter_opsinitiated                    | opsinitiated                    | Number of operations the server has initiated since it started  |
| ds_exporter_opscompleted                    | opscompleted                    | Number of operations the server has completed since it started  |
| ds_exporter_entriessent                     | entriessent                     | Number of entries sent to clients since the server started  |
| ds_exporter_bytessent                       | bytessent                       | Number of bytes sent to clients after the server starts  |
| ds_exporter_nbackends                       | nbackends                       | Number of back ends (databases) the server services  |
| ds_exporter_currenttime                     | currenttime                     | Current time of the server. The time is displayed in Greenwich Mean Time (GMT) in UTC format  |
| ds_exporter_starttime                       | starttime                       | Time when the server started. The time is displayed in Greenwich Mean Time (GMT) in UTC format  |

## cn=snmp,cn=monitor

Metrics collected from `cn=snmp,cn=monitor`

| Metric                            | Attribute             | Description |
| --------------------------------- | --------------------- | ----------- |
| ds_exporter_anonymousbinds        | anonymousbinds        | Number of anonymous bind requests    |
| ds_exporter_unauthbinds           | unauthbinds           | Number of unauthenticated (anonymous) binds    |
| ds_exporter_simpleauthbinds       | simpleauthbinds       | Number of LDAP simple bind requests (DN and password)    |
| ds_exporter_strongauthbinds       | strongauthbinds       | Number of LDAP SASL bind requests for all SASL mechanisms    |
| ds_exporter_bindsecurityerrors    | bindsecurityerrors    | Number of number of times an invalid password was given in a bind request.    |
| ds_exporter_compareops            | compareops            | Number of LDAP compare requests    |
| ds_exporter_addentryops           | addentryops           | Number of LDAP add requests.    |
| ds_exporter_removeentryops        | removeentryops        | Number of LDAP delete requests    |
| ds_exporter_modifyentryops        | modifyentryops        | Number of LDAP modify requests    |
| ds_exporter_modifyrdnops          | modifyrdnops          | Number of LDAP modify RDN (modrdn) requests    |
| ds_exporter_searchops             | searchops             | Number of LDAP search requests    |
| ds_exporter_onelevelsearchops     | onelevelsearchops     | Number of one-level search operations    |
| ds_exporter_wholesubtreesearchops | wholesubtreesearchops | Number of subtree-level search operations    |
| ds_exporter_securityerrors        | securityerrors        | Number of errors returned that were security related such as invalid passwords unknown or invalid authentication methods etc    |
| ds_exporter_errors                | errors                | Number of errors returned    |

## cn=monitor,cn=ldbm database,cn=plugins,cn=config

Metrics collected for Berkeley DB:

| Metric                           | Attribute                     | Description |
| -------------------------------- | ----------------------------- | ----------- |
| ds_dbcachehits                   | dbcachehits                   | Number of requested pages found in the database |
| ds_dbcachetries                  | dbcachetries                  | Total number of cache lookups |
| ds_dbcachehitratio               | dbcachehitratio               | Percentage of requested pages found in the database cache |
| ds_dbcachepagein                 | dbcachepagein                 | Number of pages read into the database cache |
| ds_dbcachepageout                | dbcachepageout                | Number of pages written from the database cache to the backing file |
| ds_dbcacheroevict                | dbcacheroevict                | Number of clean pages forced from the cache |
| ds_dbcacherwevict                | dbcacherwevict                | Number   of dirty pages forced from the cache |
| ds_normalizeddncachetries        | normalizeddncachetries        | Total number of cache lookups since the instance was started |
| ds_normalizeddncachehits         | normalizeddncachehits         | Normalized DNs found within the cache. |
| ds_normalizeddncachemisses       | normalizeddncachemisses       | Normalized DNs not found within the cache |
| ds_normalizeddncachehitratio     | normalizeddncachehitratio     | Percentage of the normalized DNs found in the cache |
| ds_currentnormalizeddncachesize  | currentnormalizeddncachesize  | Current size of the normalized DN cache in bytes |
| ds_maxnormalizeddncachesize      | maxnormalizeddncachesize      | Maximum size of NDn cache |
| ds_currentnormalizeddncachecount | currentnormalizeddncachecount | Number of normalized cached DNs |


Metrics collected for LMDB:

| Metric                           | Attribute                     | Description |
| -------------------------------- | ----------------------------- | ----------- |
| ds_normalizeddncachetries        | normalizeddncachetries        | Total number of cache lookups since the instance was started |
| ds_normalizeddncachehits         | normalizeddncachehits         | Normalized DNs found within the cache. |
| ds_normalizeddncachemisses       | normalizeddncachemisses       | Normalized DNs not found within the cache |
| ds_normalizeddncachehitratio     | normalizeddncachehitratio     | Percentage of the normalized DNs found in the cache |
| ds_currentnormalizeddncachesize  | currentnormalizeddncachesize  | Current size of the normalized DN cache in bytes |
| ds_maxnormalizeddncachesize      | maxnormalizeddncachesize      | Maximum size of NDn cache |
| ds_currentnormalizeddncachecount | currentnormalizeddncachecount | Number of normalized cached DNs |


## cn=monitor,cn=<backend>cn=ldbm database,cn=plugins,cn=config

Metrics collected for each backend database:

| Metric                             | Attribute              | Description |
| ---------------------------------- | ---------------------- | ----------- |
| ds_exporter_dncachehits            | dncachehits            | Number of times the server could process a request by obtaining a normalized distinguished name (DN) from the DN cache rather than normalizing it again
| ds_exporter_dncachetries           | dncachetries           | Total number of DN cache accesses since you started the instance
| ds_exporter_dncachehitratio        | dncachehitratio        | Ratio of cache tries to successful DN cache hits. The closer this value is to 100% the better
| ds_exporter_currentdncachesize     | currentdncachesize     | Total size in bytes of DN currently present in the DN cache
| ds_exporter_maxdncachesize         | maxdncachesize         | Maximum size in bytes of DNs that DS can maintain in the DN cache
| ds_exporter_currentdncachecount    | currentdncachecount    | Number of DNs currently present in the DN cache
| ds_exporter_entrycachehits         | entrycachehits         | Total number of successful entry cache lookups
| ds_exporter_entrycachetries        | entrycachetries        | Total number of entry cache lookups since you started the instance
| ds_exporter_entrycachehitratio     | entrycachehitratio     | Number of entry cache tries to successful entry cache lookups
| ds_exporter_maxentrycachesize      | maxentrycachesize      | Maximum size in bytes of directory entries that DS can maintain in the entry cache
| ds_exporter_currententrycachesize  | currententrycachesize  | Total size in bytes of directory entries currently present in the entry cache
| ds_exporter_currententrycachecount | currententrycachecount | Current number of entries stored in the entry cache of a given backend



# cn=monitor,cn=ldbm database,cn=plugins,cn=config

Metrics collected for Berkeley DB:

| Metric                                    | Attribute                             | Description |
| ----------------------------------------- | ------------------------------------- | ----------- |
| ds_exporter_db_abort_rate                 | nsslapd-db-abort-rate                 | Number of transactions that have been aborted  |
| ds_exporter_db_active_txns                | nsslapd-db-active-txns                | Number of transactions that are currently active (used by the database)  |
| ds_exporter_db_cache_size_bytes           | nsslapd-db-cache-size-bytes           | Total cache size in bytes  |
| ds_exporter_db_cache_region_wait_rate     | nsslapd-db-cache-region-wait-rate     | Number of times that a thread of control was forced to wait before obtaining the region lock  |
| ds_exporter_db_clean_pages                | nsslapd-db-clean-pages                | Clean pages currently in the cache  |
| ds_exporter_db_commits_rate               | nsslapd-db-commit-rate                | Number of transactions that have been committed  |
| ds_exporter_db_deadlock_rate              | nsslapd-db-deadlock-rate              | Number of deadlocks detected  |
| ds_exporter_db_dirty_pages                | nsslapd-db-dirty-pages                | Dirty pages currently in the cache  |
| ds_exporter_db_hash_buckets               | nsslapd-db-hash-buckets               | Number of hash buckets in buffer hash table  |
| ds_exporter_db_hash_elements_examine_rate | nsslapd-db-hash-elements-examine-rate | Total number of hash elements traversed during hash table lookups  |
| ds_exporter_db_hash_search_rate           | nsslapd-db-hash-search-rate           | Total number of buffer hash table lookups  |
| ds_exporter_db_lock_conflicts             | nsslapd-db-lock-conflicts             | Total number of locks not immediately available due to conflicts  |
| ds_exporter_db_lock_region_wait_rate      | nsslapd-db-lock-region-wait-rate      | Number of times that a thread of control was forced to wait before obtaining the region lock  |
| ds_exporter_db_lock_request_rate          | nsslapd-db-lock-request-rate          | Total number of locks requested  |
| ds_exporter_db_lockers                    | nsslapd-db-lockers                    | Number of current lockers  |
| ds_exporter_db_configured_locks           | nsslapd-db-configured-locks           | Configured number of locks  |
| ds_exporter_db_current_locks              | nsslapd-db-current-locks              | Number of locks currently used by the database  |
| ds_exporter_db_max_locks                  | nsslapd-db-max-locks                  | The maximum number of locks at any one time  |
| ds_exporter_db_log_region_wait_rate       | nsslapd-db-log-region-wait-rate       | Number of times that a thread of control was forced to wait before obtaining the region lock  |
| ds_exporter_db_log_write_rate             | nsslapd-db-log-write-rate             | Number of bytes written to the log since the last checkpoint  |
| ds_exporter_db_longest_chain_length       | nsslapd-db-longest-chain-length       | Longest chain ever encountered in buffer hash table lookups  |
| ds_exporter_db_page_create_rate           | nsslapd-db-page-create-rate           | Pages created in the cache  |
| ds_exporter_db_page_read_rate             | nsslapd-db-page-read-rate             | Pages read into the cache  |
| ds_exporter_db_page_ro_evict_rate         | nsslapd-db-page-ro-evict-rate         | Clean pages forced from the cache  |
| ds_exporter_db_page_rw_evict_rate         | nsslapd-db-page-rw-evict-rate         | Dirty pages forced from the cache  |
| ds_exporter_db_page_trickle_rate          | nsslapd-db-page-trickle-rate          | Dirty pages written using the memp_trickle interface  |
| ds_exporter_db_page_write_rate            | nsslapd-db-page-write-rate            | Pages read into the cache  |
| ds_exporter_db_pages_in_use               | nsslapd-db-pages-in-use               | All pages, clean or dirty, currently in use  |
| ds_exporter_db_txn_region_wait_rate       | nsslapd-db-txn-region-wait-rate       | Number of times that a thread of control was force to wait before obtaining the region lock  |
| ds_exporter_db_current_lock_objects       | nsslapd-db-current-lock-objects       | The number of current lock objects  |
| ds_exporter_db_max_lock_objects           | nsslapd-db-max-lock-objects           | The maximum number of lock objects at any one time  |

Metrics collected for LMDB:

| Metric                      | Attribute       | Description |
| --------------------------- | --------------- | ----------- |
| ds_exporter_dbenvmapsize    | dbenvmapsize    | LMDB Size of the data file in bytes |
| ds_exporter_dbenvlastpageno | dbenvlastpageno | LMDB Number of pages used |
| ds_exporter_dbenvlasttxnid  | dbenvlasttxnid  | LMDB Last transaction ID |
| ds_exporter_dbenvmaxreaders | dbenvmaxreaders | LMDB Max readers |
| ds_exporter_dbenvnumreaders | dbenvnumreaders | LMDB Number of readers used |
| ds_exporter_dbenvnumdbis    | dbenvnumdbis    | Number of DBIs (named databases) within the LMDB environment |
| ds_exporter_waitingrwtxn    | waitingrwtxn    | LMDB Waiting RW transactions |
| ds_exporter_activerwtxn     | activerwtxn     | LMDB Active RW transactions |
| ds_exporter_abortrwtxn      | abortrwtxn      | LMDB Aborted RW transactions |
| ds_exporter_commitrwtxn     | commitrwtxn     | LMDB Commited RW transactions |
| ds_exporter_granttimerwtxn  | granttimerwtxn  | There is no clear explanation for this parameter yet |
| ds_exporter_lifetimerwtxn   | lifetimerwtxn   | There is no clear explanation for this parameter yet |
| ds_exporter_waitingrotxn    | waitingrotxn    | MDB Waiting RO transactions |
| ds_exporter_activerotxn     | activerotxn     | LMDB Active RO transactions |
| ds_exporter_abortrotxn      | abortrotxn      | LMDB Aborted RO transactions |
| ds_exporter_commitrotxn     | commitrotxn     | LMDB Commited RO transactions |
| ds_exporter_granttimerotxn  | granttimerotxn  | There is no clear explanation for this parameter yet |
| ds_exporter_lifetimerotxn   | lifetimerotxn   | There is no clear explanation for this parameter yet |
