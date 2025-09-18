# CHANGELOG.md

## 0.2.0 (18.09.2025)

Features:

- Added new metrics for ldbm database
- Removed inefficient parameters (excessive timeouts and attempts)
- Added `/up` and `/health` endpoints to the exporter API
- Improved stability of the LDAP connection pool

Fixes:

- Fixed LDAP connection closure – previously used `Close`, now correctly uses `Unbind`


## 0.1.0 (16.09.2025)

Features:
- Added the first working version of the exporter
Fixes:

Security:
