# 389-ds-exporter API

The exporter implements the following API endpoints:

## /metrics

The main endpoint of the exporter. Returns metrics in Prometheus format.
> `/metrics` is the default path, but it can be changed in the configuration. For more details, see [config.md](config.md).

## /up

A simple endpoint to check the availability of the exporter.
It does not provide any information about connection status or internal state — it simply responds to requests if the application is running.

## /health

An endpoint to check the readiness of the exporter.
When accessed, the exporter performs an LDAP connection check and returns a JSON-formatted response.

If the connection is successful:
```json
{
  "status": {
    "ldap": "ok"
  },
  "timestamp": "2025-09-18T12:52:03Z",
  "uptime_seconds": 1893
}
```

If the connection fails:
```json
{
  "status": {
    "ldap": "unavailable"
  },
  "timestamp": "2025-09-18T12:52:03Z",
  "uptime_seconds": 1893
}
```