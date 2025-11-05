# CHANGELOG.md

## v2.0.1 (05.11.2025)

### Fixes
- Fixed set of metrics used by the dashboard (31)
- Fixed health check in containers in examples
- Added unit tests for configuration modules

## v2.0.0 (01.11.2025)

### Features
- Added support for HTTPS and Basic Auth
- Updated metric names according to Prometheus recommendations
- Added the ability to skip TLS verification when using LDAPS
- Added automatic closing of LDAP connections when idle
- Added error information output in the `/health` endpoint

### Fixes
- Fixed concurrency issues in the LDAP pool that could cause exporter failures
- Fixed minor issues in the dashboard

### Security
- Upgraded Go to version 1.25.3 to address the vulnerability [GO-2025-4007](https://pkg.go.dev/vuln/GO-2025-4007)

## v1.0.0 (02.10.2025)

### Features
- Improve stability of ldap-pool
- Add advanced logging
- Automatically detect backend database type and its instances
- Simplify configuration file format
- Update dashboard

### Fixes
- Fix mismatch between metric types and their actual values
- Fix errors when using empty sections of the configuration file

## 0.2.0 (18.09.2025)

### Features:

- Added new metrics for ldbm database
- Removed inefficient parameters (excessive timeouts and attempts)
- Added `/up` and `/health` endpoints to the exporter API
- Improved stability of the LDAP connection pool

### Fixes:

- Fixed LDAP connection closure – previously used `Close`, now correctly uses `Unbind`


## 0.1.0 (16.09.2025)

### Features:
- Added the first working version of the exporter

### Fixes:

### Security:
