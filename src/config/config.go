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

// ExporterConfiguration represents exporter configuration top level struct
type ExporterConfiguration struct {
	Global GlobalConfig `yaml:"global"`
	HTTP   HTTPConfig   `yaml:"http"`
	LDAP   LDAPConfig   `yaml:"ldap"`
}

type GlobalConfig struct {
	BackendImplement BackendType `yaml:"ds_backend_implement"`
	Backends         []string    `yaml:"ds_backends"`
}

type HTTPConfig struct {
	ListenAddress *string `yaml:"listen_address"` // Don't use these field directly. Access them via the Get function
	MetricsPath   *string `yaml:"metrics_path"`   // Don't use these field directly. Access them via the Get function
	ReadTimeout   *int    `yaml:"read_timeout"`   // Don't use these field directly. Access them via the Get function
	WriteTimeout  *int    `yaml:"write_timeout"`  // Don't use these field directly. Access them via the Get function
	IdleTimeout   *int    `yaml:"idle_timeout"`   // Don't use these field directly. Access them via the Get function
}

type LDAPConfig struct {
	ServerURL      string                   `yaml:"server_url"`
	BindDN         string                   `yaml:"bind_dn"`
	BindPw         string                   `yaml:"bind_pw"`
	ConnectionPool LDAPConnectionPoolConfig `yaml:"connection_pool"`
}

type LDAPConnectionPoolConfig struct {
	ConnectionsLimit       *uint `yaml:"connections_limit"`        // Don't use these field directly. Access them via the Get function
	DialTimeout            *uint `yaml:"dial_timeout"`             // Don't use these field directly. Access them via the Get function
	RetryCount             *uint `yaml:"retry_count"`              // Don't use these field directly. Access them via the Get function
	RetryDelay             *uint `yaml:"retry_delay"`              // Don't use these field directly. Access them via the Get function
	ConnectionAliveTimeout *uint `yaml:"connection_alive_timeout"` // Don't use these field directly. Access them via the Get function
}

// String returns string describing config with default values
func (c *ExporterConfiguration) String() string {

	var b strings.Builder
	b.WriteString(fmt.Sprintf("global.ds_backend_implemet: %s\n", c.Global.BackendImplement))
	b.WriteString(fmt.Sprintf("global.ds_backends: %v\n", c.Global.Backends))
	b.WriteString(fmt.Sprintf("http.listen_address: %s\n", c.HTTP.GetListenAddress()))
	b.WriteString(fmt.Sprintf("http.metrics_path: %s\n", c.HTTP.GetMetricsPath()))
	b.WriteString(fmt.Sprintf("http.read_timeout: %d\n", c.HTTP.GetReadTimeout()))
	b.WriteString(fmt.Sprintf("http.write_timeout: %d\n", c.HTTP.GetWriteTimeout()))
	b.WriteString(fmt.Sprintf("http.idle_timeout: %d\n", c.HTTP.GetIdleTimeout()))
	b.WriteString(fmt.Sprintf("ldap.server_url: %s\n", c.LDAP.ServerURL))
	b.WriteString(fmt.Sprintf("ldap.bind_dn: %s\n", c.LDAP.BindDN))
	b.WriteString(fmt.Sprintf("ldap.bind_pw: %s\n", "*****"))
	b.WriteString(fmt.Sprintf("ldap.connection_pool.connections_limit: %d\n", c.LDAP.ConnectionPool.GetConnectionsLimit()))
	b.WriteString(fmt.Sprintf("ldap.connection_pool.dial_timeout: %d\n", c.LDAP.ConnectionPool.GetDialTimeout()))
	b.WriteString(fmt.Sprintf("ldap.connection_pool.retry_count: %d\n", c.LDAP.ConnectionPool.GetRetryCount()))
	b.WriteString(fmt.Sprintf("ldap.connection_pool.retry_delay: %d\n", c.LDAP.ConnectionPool.GetRetryDelay()))
	b.WriteString(fmt.Sprintf("ldap.connection_pool.connection_alive_timeout: %d\n", c.LDAP.ConnectionPool.GetConnectionAliveTimeout()))

	return b.String()
}

// GetConnectionsLimit func return LDAPConnectionPoolConfig.connectionsLimit if it defined.
// Else returns config.DefaultPoolConnectionsLimit constant.
func (c *LDAPConnectionPoolConfig) GetConnectionsLimit() uint {

	if c.ConnectionsLimit == nil {
		return DefaultLDAPPoolConnectionsLimit
	}
	return *c.ConnectionsLimit
}

// GetDialTimeout func return LDAPConnectionPoolConfig.dialTimeout if it defined.
// Else returns config.DefaultPoolDialTimeout constant.
func (c *LDAPConnectionPoolConfig) GetDialTimeout() uint {

	if c.DialTimeout == nil {
		return DefaultLDAPPoolDialTimeout
	}
	return *c.DialTimeout
}

// GetRetryCount func return LDAPConnectionPoolConfig.retryCount if it defined.
// Else returns config.DefaultPoolRetryCount constant.
func (c *LDAPConnectionPoolConfig) GetRetryCount() uint {

	if c.RetryCount == nil {
		return uint(DefaultLDAPPoolRetryCount)
	}
	return *c.RetryCount
}

// GetRetryDelay func return LDAPConnectionPoolConfig.retryDelay if it defined.
// Else returns config.DefaultPoolRetryDelay constant.
func (c *LDAPConnectionPoolConfig) GetRetryDelay() uint {

	if c.RetryDelay == nil {
		return DefaultLDAPPoolRetryDelay
	}
	return *c.RetryDelay
}

// GetConnectionAliveTimeout func return LDAPConnectionPoolConfig.connectionAliveTimeout if it defined.
// Else returns config.DefaultLDAPPoolConnectionAliveTimeout constant.
func (c *LDAPConnectionPoolConfig) GetConnectionAliveTimeout() uint {

	if c.ConnectionAliveTimeout == nil {
		return DefaultLDAPPoolConnectionAliveTimeout
	}
	return *c.ConnectionAliveTimeout
}

// GetListenAddress func return HTTPConfig.listenAddress if it defined.
// Else returns config.DefaultHTTPListenAdderss constant.
func (c *HTTPConfig) GetListenAddress() string {
	if c.ListenAddress == nil {
		return DefaultHTTPListenAdderss
	}
	return *c.ListenAddress
}

// GetMetricsPath func return HTTPConfig.metricsPath if it defined.
// Else returns config.DefaultHTTPMetricsPath constant.
func (c *HTTPConfig) GetMetricsPath() string {
	if c.MetricsPath == nil {
		return DefaultHTTPMetricsPath
	}
	return *c.MetricsPath
}

// GetReadTimeout func return HTTPConfig.readTimeout if it defined.
// Else returns config.DefaultHTTPReadTimeout constant.
func (c *HTTPConfig) GetReadTimeout() int {
	if c.ReadTimeout == nil {
		return DefaultHTTPReadTimeout
	}
	return *c.ReadTimeout
}

// GetWriteTimeout func return HTTPConfig.writeTimeout if it defined.
// Else returns config.DefaultHTTPWriteTimeout constant.
func (c *HTTPConfig) GetWriteTimeout() int {
	if c.WriteTimeout == nil {
		return DefaultHTTPWriteTimeout
	}
	return *c.WriteTimeout
}

// GetIdleTimeout func return HTTPConfig.idleTimeout if it defined.
// Else returns config.DefaultHTTPIdleTimeout constant.
func (c *HTTPConfig) GetIdleTimeout() int {
	if c.IdleTimeout == nil {
		return DefaultHTTPIdleTimeout
	}
	return *c.IdleTimeout
}
