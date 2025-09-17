package config

import (
	"fmt"
	"strings"
)

type BackendType string

const (
	BackendBDB BackendType = "bdb"
	BackendMDB BackendType = "mdb"
)

// ExporterConfiguration represents exporter configuration top level struct.
type ExporterConfiguration struct {
	Global  GlobalConfig  `yaml:"global"`
	HTTP    HTTPConfig    `yaml:"http"`
	LDAP    LDAPConfig    `yaml:"ldap"`
	Logging LoggingConfig `yaml:"log"`
}

type GlobalConfig struct {
	BackendImplement       BackendType `yaml:"ds_backend_implement"`
	Backends               []string    `yaml:"ds_backends"`
	ShutdownTimeout        *int        `yaml:"shutdown_timeout"`
	NumSubordinatesRecords []string    `yaml:"ds_numsubordinate_records"`
}

type HTTPConfig struct {
	ListenAddress      *string `yaml:"listen_address"`       // Don't use these field directly.
	MetricsPath        *string `yaml:"metrics_path"`         // Don't use these field directly.
	ReadTimeout        *int    `yaml:"read_timeout"`         // Don't use these field directly.
	WriteTimeout       *int    `yaml:"write_timeout"`        // Don't use these field directly.
	IdleTimeout        *int    `yaml:"idle_timeout"`         // Don't use these field directly.
	InitialReadTimeout *int    `yaml:"initial_read_timeout"` // Don't use these field directly.
}

type LDAPConfig struct {
	ServerURL     string `yaml:"server_url"`
	BindDN        string `yaml:"bind_dn"`
	BindPw        string `yaml:"bind_pw"`
	PoolConnLimit *int   `yaml:"pool_conn_limit"`
}

type LoggingConfig struct {
	Level        *string `yaml:"level"`         // Don't use these field directly.
	Handler      *string `yaml:"handler"`       // Don't use these field directly.
	File         *string `yaml:"file"`          // Don't use these field directly.
	StdoutFormat *string `yaml:"stdout_foramt"` // Don't use these field directly.
	FileFormat   *string `yaml:"file_format"`   // Don't use these field directly.
}

// String returns string describing config with default values.
func (c *ExporterConfiguration) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("global.ds_backend_implement: %s\n", c.Global.BackendImplement))
	b.WriteString(fmt.Sprintf("global.ds_backends: %v\n", c.Global.Backends))
	b.WriteString(fmt.Sprintf("global.shutdown_timeout: %v\n", c.Global.GetShutdownTimeout()))
	b.WriteString(fmt.Sprintf("global.ds_numsubordinate_records: %v\n", c.Global.NumSubordinatesRecords))
	b.WriteString(fmt.Sprintf("http.listen_address: %s\n", c.HTTP.GetListenAddress()))
	b.WriteString(fmt.Sprintf("http.metrics_path: %s\n", c.HTTP.GetMetricsPath()))
	b.WriteString(fmt.Sprintf("http.read_timeout: %d\n", c.HTTP.GetReadTimeout()))
	b.WriteString(fmt.Sprintf("http.write_timeout: %d\n", c.HTTP.GetWriteTimeout()))
	b.WriteString(fmt.Sprintf("http.idle_timeout: %d\n", c.HTTP.GetIdleTimeout()))
	b.WriteString(fmt.Sprintf("http.initial_read_timeout: %d\n", c.HTTP.GetInitialReadTimeout()))
	b.WriteString(fmt.Sprintf("ldap.server_url: %s\n", c.LDAP.ServerURL))
	b.WriteString(fmt.Sprintf("ldap.bind_dn: %s\n", c.LDAP.BindDN))
	b.WriteString(fmt.Sprintf("ldap.bind_pw: %s\n", "*****"))
	b.WriteString(fmt.Sprintf("ldap.pool_conn_limit: %d\n", c.LDAP.PoolConnLimit))
	b.WriteString(fmt.Sprintf("log.level: %s\n", c.Logging.GetLevel()))
	b.WriteString(fmt.Sprintf("log.handler: %s\n", c.Logging.GetHandler()))
	b.WriteString(fmt.Sprintf("log.file: %s\n", c.Logging.GetFile()))
	b.WriteString(fmt.Sprintf("log.stdout_foramt: %s\n", c.Logging.GetStdoutFormat()))
	b.WriteString(fmt.Sprintf("log.file_format: %s\n", c.Logging.GetFileFormat()))

	return b.String()
}

// GetShutdownTimeout func return GlobalConfig.ShutdownTimeout if it defined.
// Else returns config.DefaultGlobalShutdownTimeout constant.
func (c *GlobalConfig) GetShutdownTimeout() int {
	if c.ShutdownTimeout == nil {
		return DefaultGlobalShutdownTimeout
	}

	return *c.ShutdownTimeout
}

// GetPoolConnLimit func return LDAPConfig.GetPoolConnLimit if it defined.
// Else returns config.DefaultLDAPPoolConnLimit constant.
func (c *LDAPConfig) GetPoolConnLimit() int {
	if c.PoolConnLimit == nil {
		return DefaultLDAPPoolConnLimit
	}

	return *c.PoolConnLimit
}

// GetListenAddress func return HTTPConfig.ListenAddress if it defined.
// Else returns config.DefaultHTTPListenAdderss constant.
func (c *HTTPConfig) GetListenAddress() string {
	if c.ListenAddress == nil {
		return DefaultHTTPListenAdderss
	}

	return *c.ListenAddress
}

// GetMetricsPath func return HTTPConfig.MetricsPath if it defined.
// Else returns config.DefaultHTTPMetricsPath constant.
func (c *HTTPConfig) GetMetricsPath() string {
	if c.MetricsPath == nil {
		return DefaultHTTPMetricsPath
	}

	return *c.MetricsPath
}

// GetReadTimeout func return HTTPConfig.ReadTimeout if it defined.
// Else returns config.DefaultHTTPReadTimeout constant.
func (c *HTTPConfig) GetReadTimeout() int {
	if c.ReadTimeout == nil {
		return DefaultHTTPReadTimeout
	}

	return *c.ReadTimeout
}

// GetWriteTimeout func return HTTPConfig.WriteTimeout if it defined.
// Else returns config.DefaultHTTPWriteTimeout constant.
func (c *HTTPConfig) GetWriteTimeout() int {
	if c.WriteTimeout == nil {
		return DefaultHTTPWriteTimeout
	}

	return *c.WriteTimeout
}

// GetIdleTimeout func return HTTPConfig.IdleTimeout if it defined.
// Else returns config.DefaultHTTPIdleTimeout constant.
func (c *HTTPConfig) GetIdleTimeout() int {
	if c.IdleTimeout == nil {
		return DefaultHTTPIdleTimeout
	}

	return *c.IdleTimeout
}

// GetInitialReadTimeout func return HTTPConfig.InitialReadTimeout if it defined.
// Else returns config.DefaultHTTPInitialReadTimeout constant.
func (c *HTTPConfig) GetInitialReadTimeout() int {
	if c.InitialReadTimeout == nil {
		return DefaultHTTPInitialReadTimeout
	}

	return *c.InitialReadTimeout
}

// GetLevel func return LoggingConfig.Level if it defined.
// Else returns config.DefaultLogLevel constant.
func (c *LoggingConfig) GetLevel() string {
	if c.Level == nil {
		return DefaultLogLevel
	}

	return *c.Level
}

// GetHandler func return LoggingConfig.Handler if it defined.
// Else returns config.DefaultLogHandler constant.
func (c *LoggingConfig) GetHandler() string {
	if c.Handler == nil {
		return DefaultLogHandler
	}

	return *c.Handler
}

// GetFile func return LoggingConfig.File if it defined.
// Else returns config.DefaultLogFile constant.
func (c *LoggingConfig) GetFile() string {
	if c.File == nil {
		return DefaultLogFile
	}

	return *c.File
}

// GetStdoutFormat func return LoggingConfig.StdoutFormat if it defined.
// Else returns config.DefaultLogStdoutFormat constant.
func (c *LoggingConfig) GetStdoutFormat() string {
	if c.StdoutFormat == nil {
		return DefaultLogStdoutFormat
	}

	return *c.StdoutFormat
}

// GetFileFormat func return LoggingConfig.FileFormat if it defined.
// Else returns config.DefaultLogFileFormat constant.
func (c *LoggingConfig) GetFileFormat() string {
	if c.FileFormat == nil {
		return DefaultLogFileFormat
	}

	return *c.FileFormat
}
