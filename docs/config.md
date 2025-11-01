# Configuration file

389-ds-exporter is configured using two YAML files:
- `config.yml` — exporter configuration (LDAP connection parameters, list of collectors, etc.)
- `web-config.yml` — web configuration based on [Prometheus web.config.file](https://prometheus.io/docs/prometheus/latest/configuration/https/) (TLS, HTTP, BasicAuth parameters, etc.)

## Exporter configuration file parameters

### collectors_default
Set of collectors enabled by default.
Possible values:
- `all` — enables all available collectors;
- `none` — disables all collectors;
- `standard` — enables the standard set of collectors.

Default value: `standard`

---

### ds_numsubordinate_records
List of LDAP entries for which the number of subordinate entries (`numSubordinates`) will be collected.
Default value: `[]`

---

### collectors_enabled
List of explicitly enabled collectors.
Used to enable specific collectors when `collectors_default` is not set to `all`.

---

### shutdown_timeout
Maximum time (in seconds) the server waits for graceful shutdown of all resources during application termination.
During this period, the server stops accepting new connections, completes ongoing requests, and cleanly closes the HTTP server, LDAP connection pool, and other resources.
If the timeout is exceeded, remaining connections will be forcibly closed.
A value of `0` means graceful shutdown is skipped.

Default value: `5`

## LDAP

### ldap_server_url
LDAP server address in RFC-2255 format.
Examples: `ldap://localhost:389` or `ldaps://remote-server`

Default value: `ldap://localhost:389`

---

### ldap_bind_dn
DN of the account used for LDAP authentication.

---

### ldap_bind_pw
Password for the LDAP account.

---

### ldap_tls_skip_verify
Skip TLS certificate verification when connecting to LDAP.

Default value: `false`

---

### ldap_pool_conn_limit
Maximum size of the LDAP connection pool.
The recommended and standard value is `5`.
If set lower, metric collectors may block while waiting for connections, slowing down the exporter.

Default value: `5`

---

### ldap_pool_get_timeout
Timeout (in seconds) for obtaining a connection from the pool.
If a connection cannot be acquired within this time, an error is returned.

Default value: `5`

---

### ldap_dial_timeout
Timeout for establishing an LDAP connection.
This time does not include the BIND operation, only socket opening.

Default value: `3`

---

### ldap_pool_idle_time
The amount of time an idle connection remains in the pool.
If the connection is not used by any request during this period, it will be closed.

Default value: `600`

---

### ldap_pool_life_time
The maximum lifetime of a connection. After this time, the connection is closed.

Default value: `3600`

## WEB configuration parameters

Web configuration parameters are defined using the standard Prometheus `web-config.yml` file.
Detailed parameter descriptions are available [here](https://prometheus.io/docs/prometheus/latest/configuration/https/).
