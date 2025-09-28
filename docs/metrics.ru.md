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
Источник: `cn=monitor`</br>

### Метрики:

#### ds_server_threads
Тип: `gauge`</br>
Атрибут: `threads`</br>

Текущее количество потоков, используемых для обработки запросов.

---

#### ds_server_currentconnections
Тип: `gauge`</br>
Атрибут: `currentconnections`</br>

Количество установленных соединений.

---

#### ds_server_totalconnections
Тип: `counter`</br>
Атрибут: `totalconnections`</br>

Количество соединений, установленных сервером с момента запуска.

---

#### ds_server_currentconnectionsatmaxthreads
Тип: `gauge`</br>
Атрибут: `currentconnectionsatmaxthreads`</br>

Текущее количество соединений, использующих максимальное разрешенное количество потоков на соединение.

---

#### ds_server_maxthreadsperconnhits
Тип: `gauge`</br>
Атрибут: `maxthreadsperconnhits`</br>

Показывает количество раз подключения достигали ограничения количества потоков на соединение.

---

#### ds_server_dtablesize
Тип: `gauge`</br>
Атрибут: `dtablesize`</br>

Количество файловых дескрипторов, доступных серверу.

---

#### ds_server_readwaiters
Тип: `gauge`</br>
Атрибут: `readwaiters`</br>

Количество подключений, некоторые запросы которых находятся в состоянии ожидания и в данный момент не обслуживаются потоком на сервере.

---

#### ds_server_opsinitiated
Тип: `counter`</br>
Атрибут: `opsinitiated`</br>

Количество операций, инициированных сервером с момента запуска.

---

#### ds_server_opscompleted
Тип: `counter`</br>
Атрибут: `opscompleted`</br>

Количество операций, завершенных сервером с момента запуска.

---

#### ds_server_entriessent
Тип: `counter`</br>
Атрибут: `entriessent`</br>

Количество записей, отправленных клиентам с момента запуска.

---

#### ds_server_bytessent
Тип: `counter`</br>
Атрибут: `bytessent`</br>

Количество байт, отправленных клиентам с момента запуска.

---

#### ds_server_nbackends
Тип: `gauge`</br>
Атрибут: `nbackends`</br>

Количество бекендов (баз данных, суффиксов), обслуживаемых сервером.

---

#### ds_server_currenttime
Тип: `gauge`</br>
Атрибут: `currenttime`</br>

Текущее время сервера в часовом поясе UTC+0 в формате Unix Timestamp.

---

#### ds_server_starttime
Тип: `gauge`</br>
Атрибут: `starttime`</br>

Время запуска сервера в часовом поясе UTC+0 в формате Unix Timestamp.


## `snmp-server`
Коллектор `snmp-server` собирает расширенные метрики сервера. Дополняет собой но не заменяет коллектор `server`</br>
Источник метрик: `cn=snmp,cn=monitor`</br>

### Метрики

#### ds_snmp_server_anonymousbinds
Тип: `counter`</br>
Атрибут: `anonymousbinds`</br>

Количество анонимных (anonymous) BIND-операций с момента запуска сервера.

---

#### ds_snmp_server_unauthbinds
Тип: `counter`</br>
Атрибут: `unauthbinds`</br>

Количество неаутентифицированных (unauth) BIND-операций с момента запуска сервера.

---

#### ds_snmp_server_simpleauthbinds
Тип: `counter`</br>
Атрибут: `simpleauthbinds`</br>

Количество "простых" (simple) BIND-операций с момента запуска сервера.

---

#### ds_snmp_server_strongauthbinds
Тип: `counter`</br>
Атрибут: `strongauthbinds`</br>

Количество защищенных (strongauth) BIND-операций с момента запуска сервера.

---

#### ds_snmp_server_bindsecurityerrors
Тип: `counter`</br>
Атрибут: `bindsecurityerrors`</br>

Количество раз, когда в BIND-запросе был указан неверный пароль.

---

#### ds_snmp_server_compareops
Тип: `counter`</br>
Атрибут: `compareops`</br>

Количество LDAP `compare` запросов с момента запуска сервера.

---

#### ds_snmp_server_addentryops
Тип: `counter`</br>
Атрибут: `addentryops`</br>

Количество LDAP `add` запросов с момента запуска сервера.

---

#### ds_snmp_server_removeentryops
Тип: `counter`</br>
Атрибут: `removeentryops`</br>

Количество LDAP `delete` запросов с момента запуска сервера.

---

#### ds_snmp_server_modifyentryops
Тип: `counter`</br>
Атрибут: `modifyentryops`</br>

Количество LDAP `modify` запросов с момента запуска сервера.

---

#### ds_snmp_server_modifyrdnops
Тип: `counter`</br>
Атрибут: `modifyrdnops`</br>

Количество LDAP `modrdn` запросов с момента запуска сервера.

---

#### ds_snmp_server_searchops
Тип: `counter`</br>
Атрибут: `searchops`</br>

Количество LDAP `search` запросов с момента запуска сервера.

---

#### ds_snmp_server_onelevelsearchops
Тип: `counter`</br>
Атрибут: `onelevelsearchops`</br>

Количество `one-level search` запросов с момента запуска сервера.

---

#### ds_snmp_server_wholesubtreesearchops
Тип: `counter`</br>
Атрибут: `wholesubtreesearchops`</br>

Количество `subtree-level search` запросов с момента запуска сервера.

---

#### ds_snmp_server_securityerrors
Тип: `counter`</br>
Атрибут: `securityerrors`</br>

Количество возвращенных ошибок, связанных с безопасностью, таких как неправильные пароли, неправильные методы аутентификации или требования более высокого уровня безопасности.

---

#### ds_snmp_server_errors
Тип: `counter`</br>
Атрибут: `errors`</br>

Количество возвращенных ошибок.


## `numsubordinates`
Cобирает информацию о количестве записей в DN, указанных в конфигурации.</br>
Источник метрик: DN, указанные в конфигурации

### Метрики

#### ds_numsubordinates_count
Тип: `gauge`</br>
Атрибут: `numsubordinates`</br>

Количество дочерних записей DN.


## `ndn-cache`
Собирает информацию об исползьвании кеша нормализованных DN (Normalized DN Cache).</br>
Источник: `cn=monitor,cn=ldbm database,cn=plugins,cn=config`

### Метрики

#### ds_ldbm_normalizeddncachetries
Тип: `gauge`</br>
Атрибут: `normalizeddncachetries`</br>

Общее количество обращений к NDN-кешу с момента запуска сервера.

---

#### ds_ldbm_normalizeddncachehits
Тип: `gauge`</br>
Атрибут: `normalizeddncachehits`</br>

Количество нормализованных DN, найденных в кеше с момента запуска сервера.

---

#### ds_ldbm_normalizeddncachemisses
Тип: `gauge`</br>
Атрибут: `normalizeddncachemisses`</br>

Количество нормализованных DN, не найденных в кеше с момента запуска сервера.

---

#### ds_ldbm_normalizeddncachehitratio
Тип: `gauge`</br>
Атрибут: `normalizeddncachehitratio`</br>

Процент нормализованных DN, найденных в кеше.

---

#### ds_ldbm_currentnormalizeddncachesize
Тип: `gauge`</br>
Атрибут: `currentnormalizeddncachesize`</br>

Текущий размер NDN кеша в байтах.

---

#### ds_ldbm_maxnormalizeddncachesize
Тип: `gauge`</br>
Атрибут: `maxnormalizeddncachesize`</br>

Установленный максимальный размер NDN-кеша.

---

#### ds_ldbm_currentnormalizeddncachecount
Тип: `gauge`</br>
Атрибут: `currentnormalizeddncachecount`</br>

Количество закешированных нормализованных DN.


## `ldbm-instance`
Собирает метрики бекенд баз данных. Список баз данных получается автоматически при запуске экспортера.</br>
Источник: `cn=monitor,cn=<имя базы данных>,cn=ldbm database,cn=plugins,cn=config`

#### ds_ldbm_instance_entrycachehits
Тип: `counter`</br>
Атрибут: `entrycachehits`</br>

Общее количество успешных обращений к entry-кешу.

---

#### ds_ldbm_instance_entrycachetries
Тип: `counter`</br>
Атрибут: `entrycachetries`</br>

Общее количество попыток обращения к entry-кешу с момента запуска сервера.

---

#### ds_ldbm_instance_entrycachehitratio
Тип: `gauge`</br>
Атрибут: `entrycachehitratio`</br>

Отношение количества удачных обращений к entry-кешу к общему числу попыток.

---

#### ds_ldbm_instance_currententrycachesize
Тип: `gauge`</br>
Атрибут: `currententrycachesize`</br>

Текущий размер entry-кеша в байтах.

---

#### ds_ldbm_instance_maxentrycachesize
Тип: `gauge`</br>
Атрибут: `maxentrycachesize`</br>

Максимальный размер entry-кеша в байтах.

---

#### ds_ldbm_instance_currententrycachecount
Тип: `gauge`</br>
Атрибут: `currententrycachecount`</br>

Текущее количество записей, сохранённых в entry-кеше.

---

#### ds_ldbm_instance_dncachehits
Тип: `counter`</br>
Атрибут: `dncachehits`</br>

Количество обращений, когда запись была найдена в кеше.

---

#### ds_ldbm_instance_dncachetries
Тип: `counter`</br>
Атрибут: `dncachetries`</br>

Общее количество обращений к dn-кешу с момента запуска сервера.

---

#### ds_ldbm_instance_dncachehitratio
Тип: `gauge`</br>
Атрибут: `dncachehitratio`</br>

Отношение количества удачных обращений к dn-кешу к общему числу обращений.

---

#### ds_ldbm_instance_currentdncachesize
Тип: `gauge`</br>
Атрибут: `currentdncachesize`</br>

Текущий размер DN-кеша в байтах.

---

#### ds_ldbm_instance_maxdncachesize
Тип: `gauge`</br>
Атрибут: `maxdncachesize`</br>

Максимальный размер DN-кеша в байтах.

---

#### ds_ldbm_instance_currentdncachecount
Тип: `gauge`</br>
Атрибут: `currentdncachecount`</br>

Текущее количество записей в dn-кеше.


## `bdb-cache`
Собирает метрики кэша BerkeleyDB.</br>
Источник: `cn=monitor,cn=bdb,cn=ldbm database,cn=plugins,cn=config`

#### ds_bdb_cachehits
Тип: `counter`</br>
Атрибут: `dbcachehits`</br>

Количество страниц, которые были найдены в кэше базы данных без обращения к файлам на диске.

---

#### ds_bdb_cachetries
Тип: `counter`</br>
Атрибут: `dbcachetries`</br>

Общее количество обращений к кэшу базы данных с момента запуска сервера.

---

#### ds_bdb_cachehitratio
Тип: `gauge`</br>
Атрибут: `dbcachehitratio`</br>

Процент запросов к страницам, которые были найдены в кэше базы данных. Чем выше значение, тем эффективнее используется кэш.

---

#### ds_bdb_cachepagein
Тип: `gauge`</br>
Атрибут: `dbcachepagein`</br>

Количество страниц, загруженных в кэш базы данных с диска.

---

#### ds_bdb_cachepageout
Тип: `gauge`</br>
Атрибут: `dbcachepageout`</br>

Количество страниц, выгруженных из кэша базы данных на диск.

---

#### ds_bdb_cacheroevict
Тип: `gauge`</br>
Атрибут: `dbcacheroevict`</br>

Количество «чистых» (не требующих записи на диск) страниц, удалённых из кэша.

---

#### ds_bdb_cacherwevict
Тип: `gauge`</br>
Атрибут: `dbcacherwevict`</br>

Количество «грязных» (требующих записи на диск) страниц, удалённых из кеша.


## `bdb-internal`
Собирает внутренние метрики BerkeleyDB, связанные с транзакциями, блокировками, страницами и логом транзакций.</br>
Источник: `cn=monitor,cn=bdb,cn=ldbm database,cn=plugins,cn=config`

#### ds_bdb_abort_rate
Тип: `counter`</br>
Атрибут: `nsslapd-db-abort-rate`</br>

Количество прерванных транзакций.

---

#### ds_bdb_active_txns
Тип: `gauge`</br>
Атрибут: `nsslapd-db-active-txns`</br>

Количество транзакций, которые в данный момент активны и используются базой данных.

---

#### ds_bdb_cache_size_bytes
Тип: `gauge`</br>
Атрибут: `nsslapd-db-cache-size-bytes`</br>

Максимальный, установленный конфигурацией, размер кэша базы данных в байтах.

---

#### ds_bdb_cache_region_wait_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-cache-region-wait-rate`</br>

Количество случаев, когда потоку приходилось ждать для получения блокировки региона кэша.

---

#### ds_bdb_clean_pages
Тип: `gauge`</br>
Атрибут: `nsslapd-db-clean-pages`</br>

Количество «чистых» страниц в кэше базы данных.

---

#### ds_bdb_commit_rate
Тип: `counter`</br>
Атрибут: `nsslapd-db-commit-rate`</br>

Количество зафиксированных транзакций.

---

#### ds_bdb_deadlock_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-deadlock-rate`</br>

Суммарное количество обнаруженных дедлоков с момента запуска сервера.

---

#### ds_bdb_dirty_pages
Тип: `gauge`</br>
Атрибут: `nsslapd-db-dirty-pages`</br>

Количество «грязных» страниц в кэше базы данных.

---

#### ds_bdb_hash_buckets
Тип: `gauge`</br>
Атрибут: `nsslapd-db-hash-buckets`</br>

Количество хэш-бакетов в хеш-таблице буфера.

---

#### ds_bdb_hash_elements_examine_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-hash-elements-examine-rate`</br>

Количество хэш-элементов, просмотренных при поисках в хэш-таблице.

---

#### ds_bdb_hash_search_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-hash-search-rate`</br>

Количество поисков в таблице буферного хэша.

---

#### ds_bdb_lock_conflicts
Тип: `gauge`</br>
Атрибут: `nsslapd-db-lock-conflicts`</br>

Количество случаев, когда блокировка не могла быть выдана из-за конфликта.

---

#### ds_bdb_lock_region_wait_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-lock-region-wait-rate`</br>

Количество случаев ожидания блокировки региона.

---

#### ds_bdb_lock_request_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-lock-request-rate`</br>

Общее количество запросов на установку блокировок.

---

#### ds_bdb_lockers
Тип: `gauge`</br>
Атрибут: `nsslapd-db-lockers`</br>

Количество текущих «локеров» (субъектов, удерживающих блокировки).

---

#### ds_bdb_configured_locks
Тип: `gauge`</br>
Атрибут: `nsslapd-db-configured-locks`</br>

Сконфигурированное количество блокировок.

---

#### ds_bdb_current_locks
Тип: `gauge`</br>
Атрибут: `nsslapd-db-current-locks`</br>

Количество блокировок, используемых в данный момент.

---

#### ds_bdb_max_locks
Тип: `gauge`</br>
Атрибут: `nsslapd-db-max-locks`</br>

Максимальное количество блокировок, использованных одновременно с момента запуска сервера.

---

#### ds_bdb_log_region_wait_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-log-region-wait-rate`</br>

Количество случаев ожидания блокировки региона лога транзакций.

---

#### ds_bdb_log_write_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-log-write-rate`</br>

Количество байт, записанных в журнал с момента последнего чекпоинта лога транзакций.

---

#### ds_bdb_longest_chain_length
Тип: `gauge`</br>
Атрибут: `nsslapd-db-longest-chain-length`</br>

Максимальная длина цепочки при поиске в хэш-таблице буферов.

---

#### ds_bdb_page_create_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-page-create-rate`</br>

Количество страниц, созданных в кэше.

---

#### ds_bdb_page_read_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-page-read-rate`</br>

Количество страниц, прочитанных в кэш.

---

#### ds_bdb_page_ro_evict_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-page-ro-evict-rate`</br>

Количество «чистых» страниц, удалённых из кэша.
> Это значение дублирует метрикуу `ds_bdb_cacheroevict`. Можно было бы оставить только одну из этих метрик,
но в 389ds она зачем-то дублируется, пусть здесь это будет также.

---

#### ds_bdb_page_rw_evict_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-page-rw-evict-rate`</br>

Количество «грязных» страниц, удалённых из кэша.
> Это значение дублирует метрикуу `ds_bdb_cacherwevict`. Можно было бы оставить только одну из этих метрик,
но в 389ds она зачем-то дублируется, пусть здесь это будет также.

---

#### ds_bdb_page_trickle_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-page-trickle-rate`</br>

Количество «грязных» страниц, записанных с использованием интерфейса `memp_trickle`.

---

#### ds_bdb_page_write_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-page-write-rate`</br>

Количество страниц, записанных из кэша.

---

#### ds_bdb_pages_in_use
Тип: `gauge`</br>
Атрибут: `nsslapd-db-pages-in-use`</br>

Количество всех страниц (чистых и грязных), которые сейчас используются кешем.

---

#### ds_bdb_txn_region_wait_rate
Тип: `gauge`</br>
Атрибут: `nsslapd-db-txn-region-wait-rate`</br>

Количество случаев ожидания блокировки региона транзакций.

---

#### ds_bdb_current_lock_objects
Тип: `gauge`</br>
Атрибут: `nsslapd-db-current-lock-objects`</br>

Текущее количество объектов блокировок.

---

#### ds_bdb_max_lock_objects
Тип: `gauge`</br>
Атрибут: `nsslapd-db-max-lock-objects`</br>

Максимальное количество объектов блокировок, зарегистрированное с момента запуска сервера.


## `lmdb-internal`
Собирает внутренние метрики LMDB, связанные с транзакциями, файлами и ресурсами окружения.</br>
Источник: `cn=monitor,cn=mdb,cn=ldbm database,cn=plugins,cn=config`

#### ds_mdb_dbenvmapsize
Тип: `gauge`</br>
Атрибут: `dbenvmapsize`</br>

Размер файла данных LMDB в байтах.

---

#### ds_mdb_dbenvlastpageno
Тип: `gauge`</br>
Атрибут: `dbenvlastpageno`</br>

Количество страниц, используемых в файле базы данных LMDB.

---

#### ds_mdb_dbenvlasttxnid
Тип: `gauge`</br>
Атрибут: `dbenvlasttxnid`</br>

Идентификатор последней транзакции LMDB.

---

#### ds_mdb_dbenvmaxreaders
Тип: `gauge`</br>
Атрибут: `dbenvmaxreaders`</br>

Максимальное количество потоков-чтения (readers), разрешённых в окружении LMDB.

---

#### ds_mdb_dbenvnumreaders
Тип: `gauge`</br>
Атрибут: `dbenvnumreaders`</br>

Текущее количество используемых потоков чтения в окружении LMDB.

---

#### ds_mdb_dbenvnumdbis
Тип: `gauge`</br>
Атрибут: `dbenvnumdbis`</br>

Количество DBI (именованных баз данных), открытых в окружении LMDB.

---

#### ds_mdb_waitingrwtxn
Тип: `gauge`</br>
Атрибут: `waitingrwtxn`</br>

Количество RW (чтение/запись) транзакций, находящихся в ожидании.

---

#### ds_mdb_activerwtxn
Тип: `gauge`</br>
Атрибут: `activerwtxn`</br>

Количество активных RW (чтение/запись) транзакций.

---

#### ds_mdb_abortrwtxn
Тип: `gauge`</br>
Атрибут: `abortrwtxn`</br>

Количество прерванных RW транзакций.

---

#### ds_mdb_commitrwtxn
Тип: `gauge`</br>
Атрибут: `commitrwtxn`</br>

Количество успешно завершённых RW транзакций.

---

#### ds_mdb_granttimerwtxn
Тип: `gauge`</br>
Атрибут: `granttimerwtxn`</br>

Описание недоступно. Для данного атрибута не получилось найти описания в документации. Если вы знаете, что означает этот атрибут - пожалуйста, создайте Issue в проекте с описанием или ссылкой на документацию.

---

#### ds_mdb_lifetimerwtxn
Тип: `gauge`</br>
Атрибут: `lifetimerwtxn`</br>

Описание недоступно. Для данного атрибута не получилось найти описания в документации. Если вы знаете, что означает этот атрибут - пожалуйста, создайте Issue в проекте с описанием или ссылкой на документацию.

---

#### ds_mdb_waitingrotxn
Тип: `gauge`</br>
Атрибут: `waitingrotxn`</br>

Количество RO (только чтение) транзакций, находящихся в ожидании.

---

#### ds_mdb_activerotxn
Тип: `gauge`</br>
Атрибут: `activerotxn`</br>

Количество активных RO транзакций.

---

#### ds_mdb_abortrotxn
Тип: `gauge`</br>
Атрибут: `abortrotxn`</br>

Количество прерванных RO транзакций.

---

#### ds_mdb_commitrotxn
Тип: `gauge`</br>
Атрибут: `commitrotxn`</br>

Количество успешно завершённых RO транзакций.

---

#### ds_mdb_granttimerotxn
Тип: `gauge`</br>
Атрибут: `granttimerotxn`</br>

Описание недоступно. Для данного атрибута не получилось найти описания в документации. Если вы знаете, что означает этот атрибут - пожалуйста, создайте Issue в проекте с описанием или ссылкой на документацию.

---

#### ds_mdb_lifetimerotxn
Тип: `gauge`</br>
Атрибут: `lifetimerotxn`</br>

Описание недоступно. Для данного атрибута не получилось найти описания в документации. Если вы знаете, что означает этот атрибут - пожалуйста, создайте Issue в проекте с описанием или ссылкой на документацию.

