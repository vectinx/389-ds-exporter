package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-ldap/ldap/v3"

	expldap "389-ds-exporter/internal/ldap"
)

// HealthHttpResponse function performs exporter healcheck and returns its json result.
func HealthHttpResponse(
	pool *expldap.Pool,
	startTime time.Time,
	timeout time.Duration,
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
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		conn, err := pool.Conn(ctx)

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
		var errMsg string
		if err == nil {
			errMsg = ""
		} else {
			errMsg = err.Error()
		}
		healthResponse := map[string]any{
			"status": map[string]string{
				"ldap": ldapStatus,
			},
			"uptime_seconds": int(uptime),
			"timestamp":      time.Now().Format(time.RFC3339),
			"error":          errMsg,
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
