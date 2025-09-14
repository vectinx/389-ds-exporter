package config

const (
	DefaultHTTPListenAdderss              string = "127.0.0.1:9389"
	DefaultHTTPMetricsPath                string = "/metrics"
	DefaultHTTPReadTimeout                int    = 10
	DefaultHTTPWriteTimeout               int    = 15
	DefaultHTTPIdleTimeout                int    = 60
	DefaultLDAPPoolConnectionsLimit       uint   = 4
	DefaultLDAPPoolDialTimeout            uint   = 1
	DefaultLDAPPoolRetryCount             uint   = 1
	DefaultLDAPPoolRetryDelay             uint   = 1
	DefaultLDAPPoolConnectionAliveTimeout uint   = 1
	DefaultLogLevel                       string = "INFO"
	DefaultLogHandler                     string = "both"
	DefaultLogFile                        string = "/var/log/389-ds-exporter.log"
	DefaultLogStdoutFormat                string = "text"
	DefaultLogFileFormat                  string = "json"
)
