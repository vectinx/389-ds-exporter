package config

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"gopkg.in/yaml.v2"
)

type EnabledCollectorsType string

const (
	DefaultGlobalShutdownTimeout  int    = 5
	DefaultHTTPListenAddress      string = "127.0.0.1:9389"
	DefaultHTTPMetricsPath        string = "/metrics"
	DefaultHTTPReadTimeout        int    = 10
	DefaultHTTPWriteTimeout       int    = 15
	DefaultHTTPIdleTimeout        int    = 60
	DefaultHTTPInitialReadTimeout int    = 3
	DefaultLDAPPoolConnLimit      int    = 4
	DefaultLDAPPoolGetTimeout     int    = 5
	DefaultLDAPPoolIdleTime       int    = 300
	DefaultLDAPPoolLifeTime       int    = 3600
	DefaultLDAPDialTimeout        int    = 3
	DefaultLogLevel               string = "INFO"
	DefaultLogHandler             string = "stdout"
	DefaultLogFile                string = "/var/log/389-ds-exporter/exporter.log"
	DefaultLogStdoutFormat        string = "text"
	DefaultLogFileFormat          string = "json"
	DefaultCollectorsDefault      string = "standard"

	BackendBDB string = "bdb"
	BackendMDB string = "mdb"

	CollectorsAll      EnabledCollectorsType = "all"
	CollectorsStandard EnabledCollectorsType = "standard"
	CollectorsNone     EnabledCollectorsType = "none"
)

type ExporterConfig struct {
	// YAML tags are needed here for correct marshalling
	// of the structure when it is necessary to display the final config

	ShutdownTimeout         int      `yaml:"shutdown_timeout"`
	CollectorsDefault       string   `yaml:"collectors_default"`
	CollectorsEnabled       []string `yaml:"collectors_enabled"`
	DSNumSubordinateRecords []string `yaml:"ds_numsubordinate_records"`
	DSBackendType           string   `yaml:"ds_backend_type"`
	DSBackendDBs            []string `yaml:"ds_backend_dbs"`

	HTTPListenAddress      string `yaml:"http_listen_address"`
	HTTPMetricsPath        string `yaml:"http_metrics_path"`
	HTTPReadTimeout        int    `yaml:"http_read_timeout"`
	HTTPWriteTimeout       int    `yaml:"http_write_timeout"`
	HTTPIdleTimeout        int    `yaml:"http_idle_timeout"`
	HTTPInitialReadTimeout int    `yaml:"http_initial_read_timeout"`

	LDAPServerURL      string `yaml:"ldap_server_url"`
	LDAPBindDN         string `yaml:"ldap_bind_dn"`
	LDAPBindPw         string `yaml:"ldap_bind_pw"`
	LDAPPoolConnLimit  int    `yaml:"ldap_pool_conn_limit"`
	LDAPPoolGetTimeout int    `yaml:"ldap_pool_get_timeout"`
	LDAPPoolIdleTime   int    `yaml:"ldap_pool_idle_time"`
	LDAPPoolLifeTime   int    `yaml:"ldap_pool_life_time"`
	LDAPDialTimeout    int    `yaml:"ldap_dial_timeout"`

	LogLevel        string `yaml:"log_level"`
	LogHandler      string `yaml:"log_handler"`
	LogFile         string `yaml:"log_file"`
	LogStdoutFormat string `yaml:"log_stdout_format"`
	LogFileFormat   string `yaml:"log_file_format"`
}

type rawConfig struct {
	ShutdownTimeout         *int     `yaml:"shutdown_timeout"`
	CollectorsDefault       *string  `yaml:"collectors_default"`
	CollectorsEnabled       []string `yaml:"collectors_enabled"`
	DSNumSubordinateRecords []string `yaml:"ds_numsubordinate_records"`
	DSBackendType           *string  `yaml:"ds_backend_type"`
	DSBackendDBs            []string `yaml:"ds_backend_dbs"`

	HTTPListenAddress      *string `yaml:"http_listen_address"`
	HTTPMetricsPath        *string `yaml:"http_metrics_path"`
	HTTPReadTimeout        *int    `yaml:"http_read_timeout"`
	HTTPWriteTimeout       *int    `yaml:"http_write_timeout"`
	HTTPIdleTimeout        *int    `yaml:"http_idle_timeout"`
	HTTPInitialReadTimeout *int    `yaml:"http_initial_read_timeout"`

	LDAPServerURL      *string `yaml:"ldap_server_url"`
	LDAPBindDN         *string `yaml:"ldap_bind_dn"`
	LDAPBindPw         *string `yaml:"ldap_bind_pw"`
	LDAPPoolConnLimit  *int    `yaml:"ldap_pool_conn_limit"`
	LDAPPoolGetTimeout *int    `yaml:"ldap_pool_get_timeout"`
	LDAPPoolIdleTime   *int    `yaml:"ldap_pool_idle_time"`
	LDAPPoolLifeTime   *int    `yaml:"ldap_pool_life_time"`
	LDAPDialTimeout    *int    `yaml:"ldap_dial_timeout"`

	LogLevel        *string `yaml:"log_level"`
	LogHandler      *string `yaml:"log_handler"`
	LogFile         *string `yaml:"log_file"`
	LogStdoutFormat *string `yaml:"log_stdout_format"`
	LogFileFormat   *string `yaml:"log_file_format"`
}

//nolint:gocognit
func (r *rawConfig) toConfig() *ExporterConfig {
	cfg := &ExporterConfig{}

	// Global
	if r.ShutdownTimeout != nil {
		cfg.ShutdownTimeout = *r.ShutdownTimeout
	} else {
		cfg.ShutdownTimeout = DefaultGlobalShutdownTimeout
	}
	if r.CollectorsDefault != nil {
		cfg.CollectorsDefault = *r.CollectorsDefault
	} else {
		cfg.CollectorsDefault = DefaultCollectorsDefault
	}

	if r.DSBackendType != nil {
		cfg.DSBackendType = *r.DSBackendType
	} else {
		cfg.DSBackendType = ""
	}

	cfg.CollectorsEnabled = r.CollectorsEnabled
	cfg.DSNumSubordinateRecords = r.DSNumSubordinateRecords
	cfg.DSBackendDBs = r.DSBackendDBs

	// HTTP
	if r.HTTPListenAddress != nil {
		cfg.HTTPListenAddress = *r.HTTPListenAddress
	} else {
		cfg.HTTPListenAddress = DefaultHTTPListenAddress
	}
	if r.HTTPMetricsPath != nil {
		cfg.HTTPMetricsPath = *r.HTTPMetricsPath
	} else {
		cfg.HTTPMetricsPath = DefaultHTTPMetricsPath
	}
	if r.HTTPReadTimeout != nil {
		cfg.HTTPReadTimeout = *r.HTTPReadTimeout
	} else {
		cfg.HTTPReadTimeout = DefaultHTTPReadTimeout
	}
	if r.HTTPWriteTimeout != nil {
		cfg.HTTPWriteTimeout = *r.HTTPWriteTimeout
	} else {
		cfg.HTTPWriteTimeout = DefaultHTTPWriteTimeout
	}
	if r.HTTPIdleTimeout != nil {
		cfg.HTTPIdleTimeout = *r.HTTPIdleTimeout
	} else {
		cfg.HTTPIdleTimeout = DefaultHTTPIdleTimeout
	}
	if r.HTTPInitialReadTimeout != nil {
		cfg.HTTPInitialReadTimeout = *r.HTTPInitialReadTimeout
	} else {
		cfg.HTTPInitialReadTimeout = DefaultHTTPInitialReadTimeout
	}

	// LDAP
	if r.LDAPServerURL != nil {
		cfg.LDAPServerURL = *r.LDAPServerURL
	}
	if r.LDAPBindDN != nil {
		cfg.LDAPBindDN = *r.LDAPBindDN
	}
	if r.LDAPBindPw != nil {
		cfg.LDAPBindPw = *r.LDAPBindPw
	}
	if r.LDAPPoolConnLimit != nil {
		cfg.LDAPPoolConnLimit = *r.LDAPPoolConnLimit
	} else {
		cfg.LDAPPoolConnLimit = DefaultLDAPPoolConnLimit
	}
	if r.LDAPPoolGetTimeout != nil {
		cfg.LDAPPoolGetTimeout = *r.LDAPPoolGetTimeout
	} else {
		cfg.LDAPPoolGetTimeout = DefaultLDAPPoolGetTimeout
	}
	if r.LDAPPoolIdleTime != nil {
		cfg.LDAPPoolIdleTime = *r.LDAPPoolIdleTime
	} else {
		cfg.LDAPPoolIdleTime = DefaultLDAPPoolIdleTime
	}
	if r.LDAPPoolLifeTime != nil {
		cfg.LDAPPoolLifeTime = *r.LDAPPoolLifeTime
	} else {
		cfg.LDAPPoolLifeTime = DefaultLDAPPoolLifeTime
	}
	if r.LDAPDialTimeout != nil {
		cfg.LDAPDialTimeout = *r.LDAPDialTimeout
	} else {
		cfg.LDAPDialTimeout = DefaultLDAPDialTimeout
	}

	// Logging
	if r.LogLevel != nil {
		cfg.LogLevel = *r.LogLevel
	} else {
		cfg.LogLevel = DefaultLogLevel
	}
	if r.LogHandler != nil {
		cfg.LogHandler = *r.LogHandler
	} else {
		cfg.LogHandler = DefaultLogHandler
	}
	if r.LogFile != nil {
		cfg.LogFile = *r.LogFile
	} else {
		cfg.LogFile = DefaultLogFile
	}
	if r.LogStdoutFormat != nil {
		cfg.LogStdoutFormat = *r.LogStdoutFormat
	} else {
		cfg.LogStdoutFormat = DefaultLogStdoutFormat
	}
	if r.LogFileFormat != nil {
		cfg.LogFileFormat = *r.LogFileFormat
	} else {
		cfg.LogFileFormat = DefaultLogFileFormat
	}

	return cfg
}

// Validate function cheks if provided configuration is valid.
func (c *ExporterConfig) Validate() error {

	if c.ShutdownTimeout < 0 {
		return errors.New("shutdown_timeout should be greater than or equal to 0")
	}

	// Also allow an empty value if the user wants the backend to be automatically detected.
	if !slices.Contains([]string{BackendBDB, BackendMDB, ""}, c.DSBackendType) {
		return fmt.Errorf("invalid ds_backend_type: %s (must be 'bdb', 'mdb' or empty)", c.DSBackendType)
	}

	if c.HTTPInitialReadTimeout <= 0 {
		return errors.New("http_initial_read_timeout should be greater than 0")
	}

	if c.LDAPServerURL == "" {
		return errors.New("ldap_server_url is required")
	}

	if c.LDAPBindDN == "" {
		return errors.New("ldap_bind_dn is required")
	}

	if c.LDAPBindPw == "" {
		return errors.New("ldap_bind_pw is required")
	}

	if c.LDAPPoolConnLimit <= 0 {
		return errors.New("invalid ldap_pool_conn_limit: must be greater than 0")
	}

	if c.LDAPPoolGetTimeout <= 0 {
		return errors.New("invalid ldap_pool_get_timeout: must be greater than 0")
	}

	if c.LDAPDialTimeout <= 0 {
		return errors.New("invalid ldap_dial_timeout: must be greater than 0")
	}

	logLevels := []string{"DEBUG", "INFO", "WARNING", "ERROR"}
	if !slices.Contains(logLevels, c.LogLevel) {
		return fmt.Errorf("invalid log.level: '%s' (must be 'DEBUG', 'INFO', 'WARNING' or 'ERROR')", c.LogLevel)
	}

	logHandlers := []string{"stdout", "file", "both"}
	if !slices.Contains(logHandlers, c.LogHandler) {
		return fmt.Errorf("invalid log.handler: '%s' (must be 'stdout', 'file' or 'both')", c.LogHandler)
	}

	logFormats := []string{"text", "json"}

	if !slices.Contains(logFormats, c.LogStdoutFormat) {
		return fmt.Errorf("invalid log.stdout_foramt: '%s' (must be 'text' or 'json')", c.LogStdoutFormat)
	}
	if !slices.Contains(logFormats, c.LogFileFormat) {
		return fmt.Errorf("invalid log.file_format: '%s' (must be 'text' or 'json')", c.LogFileFormat)
	}

	return nil
}

func ReadConfig(filename string) (*ExporterConfig, error) {
	// #nosec G304: path comes from trusted config
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var raw rawConfig
	err = yaml.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}

	return raw.toConfig(), nil
}

func (c *ExporterConfig) String() string {
	safeCfg := *c
	if safeCfg.LDAPBindPw != "" {
		safeCfg.LDAPBindPw = "*****"
	}

	out, err := yaml.Marshal(&safeCfg)
	if err != nil {
		return ""
	}
	return string(out)
}
