package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-ldap/ldap/v3"

	"389-ds-exporter/src/connections"
)

// DefaultHttpResponse function generates a standard HTML response for the exporter.
func DefaultHttpResponse(metricsPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, `<html>
	<head>
		<title>389-ds-exporter</title>
	</head>
	<body>
		<p>Metrics are <a href="%s">here</a></p>
	</body>
</html>
`, html.EscapeString(metricsPath))
		if err != nil {
			slog.Error("Error writing HTTP answer", "err", err)
		}
	}
}

// HealthHttpResponse function performs exporter healcheck and returns its json result.
func HealthHttpResponse(
	pool *connections.LdapConnectionPool,
	startTime time.Time,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		ldapStatus := "ok"
		ldapAvailable := true
		// Since the pool checks connections before issuing them,
		// and gives out either a verified (live) connection or a newly established one,
		// we can assume that if the pool has issued a connection, ldap is available.
		ldapReq := ldap.NewSearchRequest(
			"",
			ldap.ScopeBaseObject,
			ldap.NeverDerefAliases,
			1, 0, false,
			"(objectClass=*)",
			[]string{"dn"},
			nil,
		)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		conn, err := pool.Get(ctx)

		if err != nil {
			slog.Warn("Healthcheck error", "err", err)
			ldapStatus = "unavailable"
			ldapAvailable = false
		} else {
			defer conn.Close()

			_, err = conn.Search(ldapReq)
			if err != nil {
				slog.Warn("LDAP health check failed", "err", err)
				ldapStatus = "unavailable"
				ldapAvailable = false
			}
		}

		uptime := time.Since(startTime).Seconds()

		healthResponse := map[string]any{
			"status": map[string]string{
				"ldap": ldapStatus,
			},
			"uptime_seconds": int(uptime),
			"timestamp":      time.Now().Format(time.RFC3339),
		}

		if ldapAvailable {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(healthResponse)
		if err != nil {
			slog.Error("Failed to write health response", "err", err)
		}
	}
}
