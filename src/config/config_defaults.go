package config

const (
	DefaultGlobalShutdownTimeout          int    = 5
	DefaultHTTPListenAdderss              string = "127.0.0.1:9389"
	DefaultHTTPMetricsPath                string = "/metrics"
	DefaultHTTPReadTimeout                int    = 10
	DefaultHTTPWriteTimeout               int    = 15
	DefaultHTTPIdleTimeout                int    = 60
	DefaultHTTPInitialReadTimeout         int    = 3
	DefaultLDAPPoolConnLimit              int    = 4
	DefaultLDAPPoolDialTimeout            int    = 1
	DefaultLDAPPoolConnectionAliveTimeout int    = 1
	DefaultLogLevel                       string = "INFO"
	DefaultLogHandler                     string = "both"
	DefaultLogFile                        string = "/var/log/389-ds-exporter/exporter.log"
	DefaultLogStdoutFormat                string = "text"
	DefaultLogFileFormat                  string = "json"
)
