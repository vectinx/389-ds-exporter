# Метрики, собираемые экспортером

Для получения метрик экспортер использует коллекторы (сборщики). *Коллектор* - это компонент экспортера, отвечающий за сбор определенных показателей сервера.
Ниже преведен список коллекторов, доступных в текущей версии 389-ds-exporter и их назначение:
- [server](#server) - собирает основные метрики сервера.
- [snmp-server](#snmp-server) - собирает расширенные метрики сервера. Дополняет собой но не заменяет коллектор `server`.
- [numsubordinates](#numsubordinates) - собирает информацию о количестве записей в DN, указанных в конфигурации ([см. config.md](config.md#global-settings)).
- [ndn-cache](#ndn-cache) - собирает информацию об использовании кеша нормализованных DN (Normalized DN Cache).
- [ldbm-instance](#ldbm-instance) - собирает метрики бекенд баз данных.
- [bdb-cache](#bdb-cache) - собирает информацию о кешах Berkeley DB.
- [bdb-internal](#bdb-internal) - собирает внутренние метрики Berkeley DB.
- [lmdb-internal](#lmdb-internal) - собирает внутренние метрики LMDB.

Ниже приведено подробное описание метрик, собираемых каждым коллектором.

## `server`
Коллектор `server` собирает основные метрики сервера</br>
Источник: `cn=monitor`

#### ds_server_version

Тип: `gauge`</br>
Атрибут: `version`

Метрика с константным значением '1' и меткой, содержащий версию 389 Directory Server.

#### ds_server_threads

Тип: `gauge`</br>
Атрибут: `threads`

Текущее количество потоков, используемых для обработки запросов.

#### ds_server_connections

Тип: `gauge`</br>
Атрибут: `currentconnections`

Количество установленных соединений.

#### ds_server_connections_total

Тип: `counter`</br>
Атрибут: `totalconnections`

Количество соединений, установленных сервером с момента запуска.

#### ds_server_connections_max_threads

Тип: `gauge`</br>
Атрибут: `currentconnectionsatmaxthreads`

Текущее количество соединений, использующих максимальное разрешенное количество потоков на соединение.

#### ds_server_max_threads_per_conn_hits_total

Тип: `gauge`</br>
Атрибут: `maxthreadsperconnhits`

Показывает количество раз подключения достигали ограничения количества потоков на соединение.

#### ds_server_dtablesize

Тип: `gauge`</br>
Атрибут: `dtablesize`

Количество файловых дескрипторов, доступных серверу.

#### ds_server_read_waiters

Тип: `gauge`</br>
Атрибут: `readwaiters`

Количество подключений, некоторые запросы которых находятся в состоянии ожидания и в данный момент не обслуживаются потоком на сервере.

#### ds_server_ops_initiated_total

Тип: `counter`</br>
Атрибут: `opsinitiated`

Количество операций, инициированных сервером с момента запуска.

#### ds_server_ops_completed_total

Тип: `counter`</br>
Атрибут: `opscompleted`

Количество операций, завершенных сервером с момента запуска.

#### ds_server_entries_sent_total

Тип: `counter`</br>
Атрибут: `entriessent`

Количество записей, отправленных клиентам с момента запуска.

#### ds_server_bytes_sent_total

Тип: `counter`</br>
Атрибут: `bytessent`

Количество байт, отправленных клиентам с момента запуска.

#### ds_server_backends

Тип: `gauge`</br>
Атрибут: `nbackends`

Количество бекендов (баз данных, суффиксов), обслуживаемых сервером.

#### ds_server_current_time_seconds

Тип: `gauge`</br>
Атрибут: `currenttime`

Текущее время сервера в часовом поясе UTC+0 в формате Unix Timestamp.

#### ds_server_start_time_seconds

Тип: `gauge`</br>
Атрибут: `starttime`

Время запуска сервера в часовом поясе UTC+0 в формате Unix Timestamp.

## `snmp-server`
Коллектор `snmp-server` собирает расширенные метрики сервера. Дополняет собой но не заменяет коллектор `server`</br>
Источник метрик: `cn=snmp,cn=monitor`

#### ds_snmp_server_bind_anonymous_total

Тип: `counter`</br>
Атрибут: `anonymousbinds`

Количество анонимных (anonymous) BIND-операций с момента запуска сервера.

#### ds_snmp_server_bind_unauth_total

Тип: `counter`</br>
Атрибут: `unauthbinds`

Количество неаутентифицированных (unauth) BIND-операций с момента запуска сервера.

#### ds_snmp_server_bind_simple_total

Тип: `counter`</br>
Атрибут: `simpleauthbinds`

Количество "простых" (simple) BIND-операций с момента запуска сервера.

#### ds_snmp_server_bind_strong_total

Тип: `counter`</br>
Атрибут: `strongauthbinds`

Количество защищенных (strongauth) BIND-операций с момента запуска сервера.

#### ds_snmp_server_bind_security_errors_total

Тип: `counter`</br>
Атрибут: `bindsecurityerrors`

Количество раз, когда в BIND-запросе был указан неверный пароль.

#### ds_snmp_server_compare_operations_total

Тип: `counter`</br>
Атрибут: `compareops`

Количество LDAP `compare` запросов с момента запуска сервера.

#### ds_snmp_server_add_operations_total

Тип: `counter`</br>
Атрибут: `addentryops`

Количество LDAP `add` запросов с момента запуска сервера.

#### ds_snmp_server_delete_operations_total

Тип: `counter`</br>
Атрибут: `removeentryops`

Количество LDAP `delete` запросов с момента запуска сервера.

#### ds_snmp_server_modify_operations_total

Тип: `counter`</br>
Атрибут: `modifyentryops`

Количество LDAP `modify` запросов с момента запуска сервера.

#### ds_snmp_server_modify_rdn_operations_total

Тип: `counter`</br>
Атрибут: `modifyrdnops`

Количество LDAP `modrdn` запросов с момента запуска сервера.

#### ds_snmp_server_search_operations_total

Тип: `counter`</br>
Атрибут: `searchops`

Количество LDAP `search` запросов с момента запуска сервера.

#### ds_snmp_server_search_onelevel_operations_total

Тип: `counter`</br>
Атрибут: `onelevelsearchops`

Количество `one-level search` запросов с момента запуска сервера.

#### ds_snmp_server_search_whole_subtree_operations_total

Тип: `counter`</br>
Атрибут: `wholesubtreesearchops`

Количество `subtree-level search` запросов с момента запуска сервера.

#### ds_snmp_server_security_errors_total

Тип: `counter`</br>
Атрибут: `securityerrors`

Количество возвращенных ошибок, связанных с безопасностью, таких как неправильные пароли, неправильные методы аутентификации или требования более высокого уровня безопасности.

#### ds_snmp_server_errors_total

Тип: `counter`</br>
Атрибут: `errors`

Количество возвращенных ошибок.

## `numsubordinates`
Cобирает информацию о количестве записей в DN, указанных в конфигурации.</br>
Источник метрик: DN, указанные в конфигурации

#### ds_numsubordinates_count

Тип: `gauge`</br>
Атрибут: `numsubordinates`

Количество дочерних записей DN.

## `ndn-cache`
Собирает информацию об исползьвании кеша нормализованных DN (Normalized DN Cache).</br>
Источник: `cn=monitor,cn=ldbm database,cn=plugins,cn=config`

#### ds_ldbm_ndn_cache_lookups_total

Тип: `counter`</br>
Атрибут: `normalizeddncachetries`

Общее количество обращений к NDN-кешу с момента запуска сервера.

#### ds_ldbm_ndn_cache_hits_total

Тип: `counter`</br>
Атрибут: `normalizeddncachehits`

Количество нормализованных DN, найденных в кеше с момента запуска сервера.

#### ds_ldbm_ndn_cache_misses_total

Тип: `counter`</br>
Атрибут: `normalizeddncachemisses`

Количество нормализованных DN, не найденных в кеше с момента запуска сервера.

#### ds_ldbm_ndn_cache_hit_ratio

Тип: `gauge`</br>
Атрибут: `normalizeddncachehitratio`

Процент нормализованных DN, найденных в кеше.

#### ds_ldbm_ndn_cache_size_bytes

Тип: `gauge`</br>
Атрибут: `currentnormalizeddncachesize`

Текущий размер NDN кеша в байтах.

#### ds_ldbm_ndn_cache_max_size_bytes

Тип: `gauge`</br>
Атрибут: `maxnormalizeddncachesize`

Установленный максимальный размер NDN-кеша.

#### ds_ldbm_ndn_cache_entries

Тип: `gauge`</br>
Атрибут: `currentnormalizeddncachecount`

Количество закешированных нормализованных DN.

## `ldbm-instance`
Собирает метрики бекенд баз данных. Список баз данных получается автоматически при запуске экспортера.</br>
Источник: `cn=monitor,cn=<имя базы данных>,cn=ldbm database,cn=plugins,cn=config`

#### ds_ldbm_instance_entry_cache_hits_total

Тип: `counter`</br>
Атрибут: `entrycachehits`

Общее количество успешных обращений к entry-кешу.

#### ds_ldbm_instance_entry_cache_lookups_total

Тип: `counter`</br>
Атрибут: `entrycachetries`

Общее количество попыток обращения к entry-кешу с момента запуска сервера.

#### ds_ldbm_instance_entry_cache_hit_ratio

Тип: `gauge`</br>
Атрибут: `entrycachehitratio`

Отношение количества удачных обращений к entry-кешу к общему числу попыток.

#### ds_ldbm_instance_entry_cache_size_bytes

Тип: `gauge`</br>
Атрибут: `currententrycachesize`

Текущий размер entry-кеша в байтах.

#### ds_ldbm_instance_entry_cache_max_size_bytes

Тип: `gauge`</br>
Атрибут: `maxentrycachesize`

Максимальный размер entry-кеша в байтах.

#### ds_ldbm_instance_entry_cache_count

Тип: `gauge`</br>
Атрибут: `currententrycachecount`

Текущее количество записей, сохранённых в entry-кеше.

#### ds_ldbm_instance_dn_cache_hits_total

Тип: `counter`</br>
Атрибут: `dncachehits`

Количество обращений, когда запись была найдена в кеше.

#### ds_ldbm_instance_dn_cache_lookups_total

Тип: `counter`</br>
Атрибут: `dncachetries`

Общее количество обращений к dn-кешу с момента запуска сервера.

#### ds_ldbm_instance_dn_cache_hit_ratio

Тип: `gauge`</br>
Атрибут: `dncachehitratio`

Отношение количества удачных обращений к dn-кешу к общему числу обращений.

#### ds_ldbm_instance_dn_cache_size_bytes

Тип: `gauge`</br>
Атрибут: `currentdncachesize`

Текущий размер DN-кеша в байтах.

#### ds_ldbm_instance_dn_cache_max_size_bytes

Тип: `gauge`</br>
Атрибут: `maxdncachesize`

Максимальный размер DN-кеша в байтах.

#### ds_ldbm_instance_dn_cache_count

Тип: `gauge`</br>
Атрибут: `currentdncachecount`

Текущее количество записей в dn-кеше.

## `bdb-cache`
Собирает метрики кэша BerkeleyDB.</br>
Источник: `cn=monitor,cn=bdb,cn=ldbm database,cn=plugins,cn=config`

#### ds_bdb_dbcache_hits_total

Тип: `counter`</br>
Атрибут: `dbcachehits`

Количество страниц, которые были найдены в кэше базы данных без обращения к файлам на диске.

#### ds_bdb_dbcache_lookups_total

Тип: `counter`</br>
Атрибут: `dbcachetries`

Общее количество обращений к кэшу базы данных с момента запуска сервера.

#### ds_bdb_dbcache_hit_ratio

Тип: `gauge`</br>
Атрибут: `dbcachehitratio`

Процент запросов к страницам, которые были найдены в кэше базы данных. Чем выше значение, тем эффективнее используется кэш.

#### ds_bdb_dbcache_pages_in_total

Тип: `gauge`</br>
Атрибут: `dbcachepagein`

Количество страниц, загруженных в кэш базы данных с диска.

#### ds_bdb_dbcache_pages_out_total

Тип: `gauge`</br>
Атрибут: `dbcachepageout`

Количество страниц, выгруженных из кэша базы данных на диск.

#### ds_bdb_dbcache_evictions_clean_total

Тип: `gauge`</br>
Атрибут: `dbcacheroevict`

Количество «чистых» (не требующих записи на диск) страниц, удалённых из кэша.

#### ds_bdb_dbcache_evictions_dirty_total

Тип: `gauge`</br>
Атрибут: `dbcacherwevict`

Количество «грязных» (требующих записи на диск) страниц, удалённых из кеша.

## `bdb-internal`
Собирает внутренние метрики BerkeleyDB, связанные с транзакциями, блокировками, страницами и логом транзакций.</br>
Источник: `cn=monitor,cn=bdb,cn=ldbm database,cn=plugins,cn=config`

#### ds_bdb_txn_abort_total

Тип: `counter`</br>
Атрибут: `nsslapd-db-abort-rate`

Количество прерванных транзакций.

#### ds_bdb_txn_active

Тип: `gauge`</br>
Атрибут: `nsslapd-db-active-txns`

Количество транзакций, которые в данный момент активны и используются базой данных.

#### ds_bdb_cache_size_bytes

Тип: `gauge`</br>
Атрибут: `nsslapd-db-cache-size-bytes`

Максимальный, установленный конфигурацией, размер кэша базы данных в байтах.

#### ds_bdb_cache_region_wait_total

Тип: `gauge`</br>
Атрибут: `nsslapd-db-cache-region-wait-rate`

Количество случаев, когда потоку приходилось ждать для получения блокировки региона кэша.

#### ds_bdb_cache_pages_clean

Тип: `gauge`</br>
Атрибут: `nsslapd-db-clean-pages`

Количество «чистых» страниц в кэше базы данных.

#### ds_bdb_txn_commit_total

Тип: `counter`</br>
Атрибут: `nsslapd-db-commit-rate`

Количество зафиксированных транзакций.

#### ds_bdb_deadlock_total

Тип: `gauge`</br>
Атрибут: `nsslapd-db-deadlock-rate`

Суммарное количество обнаруженных дедлоков с момента запуска сервера.

#### ds_bdb_cache_pages_dirty

Тип: `gauge`</br>
Атрибут: `nsslapd-db-dirty-pages`

Количество «грязных» страниц в кэше базы данных.

#### ds_bdb_cache_hash_buckets

Тип: `gauge`</br>
Атрибут: `nsslapd-db-hash-buckets`

Количество хэш-бакетов в хеш-таблице буфера.

#### ds_bdb_cache_hash_elements_examined_total

Тип: `gauge`</br>
Атрибут: `nsslapd-db-hash-elements-examine-rate`

Количество хэш-элементов, просмотренных при поисках в хэш-таблице.

#### ds_bdb_cache_hash_lookups_total

Тип: `gauge`</br>
Атрибут: `nsslapd-db-hash-search-rate`

Количество поисков в таблице буферного хэша.

#### ds_bdb_lock_conflicts_total

Тип: `gauge`</br>
Атрибут: `nsslapd-db-lock-conflicts`

Количество случаев, когда блокировка не могла быть выдана из-за конфликта.

#### ds_bdb_lock_region_wait_total

Тип: `gauge`</br>
Атрибут: `nsslapd-db-lock-region-wait-rate`

Количество случаев ожидания блокировки региона.

#### ds_bdb_lock_request_total

Тип: `gauge`</br>
Атрибут: `nsslapd-db-lock-request-rate`

Общее количество запросов на установку блокировок.

#### ds_bdb_lockers

Тип: `gauge`</br>
Атрибут: `nsslapd-db-lockers`

Количество текущих «локеров» (субъектов, удерживающих блокировки).

#### ds_bdb_locks_configured

Тип: `gauge`</br>
Атрибут: `nsslapd-db-configured-locks`

Сконфигурированное количество блокировок.

#### ds_bdb_locks_current

Тип: `gauge`</br>
Атрибут: `nsslapd-db-current-locks`

Количество блокировок, используемых в данный момент.

#### ds_bdb_locks_max

Тип: `gauge`</br>
Атрибут: `nsslapd-db-max-locks`

Максимальное количество блокировок, использованных одновременно с момента запуска сервера.

#### ds_bdb_log_region_wait_total

Тип: `gauge`</br>
Атрибут: `nsslapd-db-log-region-wait-rate`

Количество случаев ожидания блокировки региона лога транзакций.

#### ds_bdb_log_write_bytes_total

Тип: `gauge`</br>
Атрибут: `nsslapd-db-log-write-rate`

Количество байт, записанных в журнал с момента последнего чекпоинта лога транзакций.

#### ds_bdb_cache_hash_longest_chain

Тип: `gauge`</br>
Атрибут: `nsslapd-db-longest-chain-length`

Максимальная длина цепочки при поиске в хэш-таблице буферов.

#### ds_bdb_cache_page_create_total

Тип: `counter`</br>
Атрибут: `nsslapd-db-page-create-rate`

Количество страниц, созданных в кэше.

#### ds_bdb_cache_page_read_total

Тип: `counter`</br>
Атрибут: `nsslapd-db-page-read-rate`

Количество страниц, прочитанных в кэш.

#### ds_bdb_cache_page_ro_evict_total

Тип: `counter`</br>
Атрибут: `nsslapd-db-page-ro-evict-rate`

Количество «чистых» страниц, удалённых из кэша.
> Это значение дублирует метрикуу `ds_bdb_cacheroevict`. Можно было бы оставить только одну из этих метрик,
но в 389ds она зачем-то дублируется, пусть здесь это будет также.

#### ds_bdb_cache_page_rw_evict_total

Тип: `counter`</br>
Атрибут: `nsslapd-db-page-rw-evict-rate`

Количество «грязных» страниц, удалённых из кэша.
> Это значение дублирует метрикуу `ds_bdb_cacherwevict`. Можно было бы оставить только одну из этих метрик,
но в 389ds она зачем-то дублируется, пусть здесь это будет также.

#### ds_bdb_cache_page_trickle_total

Тип: `counter`</br>
Атрибут: `nsslapd-db-page-trickle-rate`

Количество «грязных» страниц, записанных с использованием интерфейса `memp_trickle`.

#### ds_bdb_cache_page_write_total

Тип: `counter`</br>
Атрибут: `nsslapd-db-page-write-rate`

Количество страниц, записанных из кэша.

#### ds_bdb_cache_pages_in_use

Тип: `gauge`</br>
Атрибут: `nsslapd-db-pages-in-use`

Количество всех страниц (чистых и грязных), которые сейчас используются кешем.

#### ds_bdb_txn_region_wait_total

Тип: `counter`</br>
Атрибут: `nsslapd-db-txn-region-wait-rate`

Количество случаев ожидания блокировки региона транзакций.

#### ds_bdb_lock_objects_current

Тип: `gauge`</br>
Атрибут: `nsslapd-db-current-lock-objects`

Текущее количество объектов блокировок.

#### ds_bdb_lock_objects_max

Тип: `gauge`</br>
Атрибут: `nsslapd-db-max-lock-objects`

Максимальное количество объектов блокировок, зарегистрированное с момента запуска сервера.

## `lmdb-internal`
Собирает внутренние метрики LMDB, связанные с транзакциями, файлами и ресурсами окружения.</br>
Источник: `cn=monitor,cn=mdb,cn=ldbm database,cn=plugins,cn=config`

#### ds_lmdb_env_map_size_bytes

Тип: `gauge`</br>
Атрибут: `dbenvmapsize`

Размер файла данных LMDB в байтах.

#### ds_lmdb_env_last_page_number

Тип: `gauge`</br>
Атрибут: `dbenvlastpageno`

Количество страниц, используемых в файле базы данных LMDB.

#### ds_lmdb_env_last_txn_id

Тип: `gauge`</br>
Атрибут: `dbenvlasttxnid`

Идентификатор последней транзакции LMDB.

#### ds_lmdb_env_max_readers

Тип: `gauge`</br>
Атрибут: `dbenvmaxreaders`

Максимальное количество потоков-чтения (readers), разрешённых в окружении LMDB.

#### ds_lmdb_env_num_readers

Тип: `gauge`</br>
Атрибут: `dbenvnumreaders`

Текущее количество используемых потоков чтения в окружении LMDB.

#### ds_lmdb_env_num_dbis

Тип: `gauge`</br>
Атрибут: `dbenvnumdbis`

Количество DBI (именованных баз данных), открытых в окружении LMDB.

#### ds_lmdb_rw_txn_waiting

Тип: `gauge`</br>
Атрибут: `waitingrwtxn`

Количество RW (чтение/запись) транзакций, находящихся в ожидании.

#### ds_lmdb_rw_txn_active

Тип: `gauge`</br>
Атрибут: `activerwtxn`

Количество активных RW (чтение/запись) транзакций.

#### ds_lmdb_rw_txn_aborted

Тип: `gauge`</br>
Атрибут: `abortrwtxn`

Количество прерванных RW транзакций.

#### ds_lmdb_rw_txn_committed

Тип: `gauge`</br>
Атрибут: `commitrwtxn`

Количество успешно завершённых RW транзакций.

#### ds_lmdb_rw_txn_grant_time

Тип: `gauge`</br>
Атрибут: `granttimerwtxn`

Описание недоступно. Для данного атрибута не получилось найти описания в документации. Если вы знаете, что означает этот атрибут - пожалуйста, создайте Issue в проекте с описанием или ссылкой на документацию.

#### ds_lmdb_rw_txn_lifetime

Тип: `gauge`</br>
Атрибут: `lifetimerwtxn`

Описание недоступно. Для данного атрибута не получилось найти описания в документации. Если вы знаете, что означает этот атрибут - пожалуйста, создайте Issue в проекте с описанием или ссылкой на документацию.

#### ds_lmdb_ro_txn_waiting

Тип: `gauge`</br>
Атрибут: `waitingrotxn`

Количество RO (только чтение) транзакций, находящихся в ожидании.

#### ds_lmdb_ro_txn_active

Тип: `gauge`</br>
Атрибут: `activerotxn`

Количество активных RO транзакций.

#### ds_lmdb_ro_txn_aborted

Тип: `gauge`</br>
Атрибут: `abortrotxn`

Количество прерванных RO транзакций.

#### ds_lmdb_ro_txn_committed

Тип: `gauge`</br>
Атрибут: `commitrotxn`</br>

Количество успешно завершённых RO транзакций.

#### ds_lmdb_ro_txn_grant_time

Тип: `gauge`</br>
Атрибут: `granttimerotxn`

Описание недоступно. Для данного атрибута не получилось найти описания в документации. Если вы знаете, что означает этот атрибут - пожалуйста, создайте Issue в проекте с описанием или ссылкой на документацию.

#### ds_lmdb_ro_txn_lifetime
Тип: `gauge`</br>
Атрибут: `lifetimerotxn`</br>

Описание недоступно. Для данного атрибута не получилось найти описания в документации. Если вы знаете, что означает этот атрибут - пожалуйста, создайте Issue в проекте с описанием или ссылкой на документацию.
