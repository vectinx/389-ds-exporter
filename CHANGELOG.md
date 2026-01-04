# CHANGELOG.md

## v2.0.4 (03.01.2026)

### Features
- The `ds_server_version` metric has been added, reflecting the 389-ds version that the exporter is working with

### Security
- Updated the go version to 1.25.5, fixing vulnerabilities GO-2025-4175 and GO-2025-4155
- Updated dependency versions


## v2.0.3 (26.11.2025)

### Security
- `golang.org/x/crypto` updated to v0.45.0 to fix CVE-2025-58181 and CVE-2025-47914 vulnerabilities

## v2.0.2 (14.11.2025)

### Features

- Upgraded go to version 1.25.4

### Fixes
- Fixed the build of docker images - added explicit file permissions
- Fixed build pipelines to meet safety requirements

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
