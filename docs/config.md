# Configuration File

The 389-ds-exporter is configured via a YAML file. This file defines global settings, HTTP server options, LDAP connection details, and logging parameters.

## Global settings

The 'global' section defines general exporter parameters.

```yaml
global:
  ds_backend_implement: bdb
  ds_backends: []
  ds_numsubordinate_records: []
  shutdown_timeout: 5
```

### ds_backend_implement
Specifies the backend database used by 389-ds.
Possible values: bdb, mbd

### ds_backends
List of 389-ds backend databases to monitor.
See: Red Hat Backend Database Docs
Default value: `[]`

### ds_numsubordinate_records
List of LDAP entries for which the number of subordinates (numSubordinates) will be collected.
Default value: `[]`

### shutdown_timeout
Timeout (in seconds) for graceful shutdown. If 0, shutdown happens immediately.
Default value: `5`

## HTTP server configuration

The 'http' section defines the parameters of the HTTP server that provides metrics.

```yaml
http:
  listen_address: ":9389"
  metrics_path: "/metrics"
  read_timeout: 10
  write_timeout: 15
  idle_timeout: 60
  initial_read_timeout: 3
```

### listen_address
Address and port to bind the HTTP server (e.g., :9389, 0.0.0.0:9389).

Default value: `127.0.0.1:9389`

### metrics_path
Path to expose Prometheus metrics (default: /metrics).

Default value: `/metrics`

### read_timeout
Maximum duration allowed for reading the entire request, including the body.
A zero or negative value means no timeout.

Default value: `10`

### write_timeout
Maximum duration before timing out writes of the response.
This timeout is reset whenever a new request’s header is read.
A zero or negative value means no timeout.

Default value: `15`

### idle_timeout
Maximum duration to wait for the next request when keep-alives are enabled.
If zero, the value of ReadTimeout is used.
If negative, or if both this and ReadTimeout are zero or negative, there is no timeout.

Default value: `60`

# initial_read_timeout
Maximum duration to wait for the client to send the beginning of the request after a connection is accepted.
If the client doesn't send any data within this time, the connection will be closed.
A value of 0 disables the timeout, meaning the server will wait indefinitely.

Default value: `3`

## LDAP Configuration

The 'ldap' section defines the parameters of the connection to the LDAP server.

```
ldap:
  server_url: "ldap://localhost:389"
  bind_dn: "cn=directory manager"
  bind_pw: "12345678"

  connection_pool:
    connections_limit: 4
    dial_timeout: 1
    retry_count: 0
    retry_delay: 1
    connection_alive_timeout: 1
```
### server_url
LDAP server URI (e.g., ldap://localhost:389 or ldaps://example.com)

Default value: `ldap://localhost:389`

### bind_dn
DN of the account used to authenticate with the LDAP server.

### bind_pw
Password of the LDAP account.

### connection_pool
The 'ldap.connection_pool' section describes the parameters for connecting to the LDAP server from which metrics will be collected.
In most cases, these settings do not need to be changed.

### connections_limit:
Maximum size of the LDAP connection pool.
Connections are created as needed and deleted when their lifetime or idle time is exceeded.
Recommended and default value is 4. If set lower, metric collectors may block waiting for connections, slowing down the exporter.
Values above 4 have no effect, as the current version of the exporter does not use more than 4 connections at once.

Default value: `4`


### dial_timeout:
Timeout in seconds for establishing a connection to the LDAP server.
Prevents the connection attempt from hanging indefinitely.
A zero or negative value means no timeout.

Default value: `1`

### retry_count:
Number of retry attempts when connecting to the LDAP server.
If the initial attempt fails, this number of retries will be made with delays between attempts.
A value of 0 means no retries.

Default value: `0`

### retry_delay:
Delay in seconds between LDAP reconnection attempts (used with retry_count).
A value of 0 means no delay between attempts.

Default value: `1`

### connection_alive_timeout:
Timeout in seconds for checking if an existing connection is still alive.
Before reusing a connection, the pool checks it by sending a basic query.
A zero or negative value disables the timeout.

Default value: `1`

## Logging
The 'log' section defines logging parameters.

### level
Logging level.
Options: `DEBUG`, `INFO`, `WARNING`, or `ERROR`.

Default value: `INFO`

### handler: both
Log output target.
Options:
- `stdout` - logs are printed to standard output
- `file`   - logs are written to a file
- `both`   - logs are written to both stdout and a file

Defeault value: `both`

### file
Log file path.
Only relevant if 'handler' is set to 'file' or 'both'.
Make sure the system is configured to rotate this file (e.g., using logrotate).

Default value: `/var/log/389-ds-exporter/exporter.log`

### stdout_format
STDOUT log format.
Options: `text`, `json`
Only applies if 'handler' is `stdout` or `both`.

Default value: `text`

### file_format
File log format.
Options: `text`, `json`
Applies only if 'handler' is `file` or `both`.

Default value: `json`