package config

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"gopkg.in/yaml.v2"
)

// Validate function cheks if provided configuration is valid.
func (c *ExporterConfiguration) Validate() error {
	if c.LDAP.ServerURL == "" {
		return errors.New("configuration parameter ldap.server_url is required")
	}

	if c.LDAP.BindDN == "" {
		return errors.New("configuration parameter ldap.bind_dn is required")
	}

	if c.LDAP.BindPw == "" {
		return errors.New("configuration parameter ldap.bind_pw is required")
	}

	if c.LDAP.PoolConnLimit <= 0 {
		return errors.New("invalid ldap.pool_conn_limit: must be greater than 0")
	}

	logLevels := []string{"DEBUG", "INFO", "WARNING", "ERROR"}
	if !slices.Contains(logLevels, c.Logging.Level) {
		return fmt.Errorf("invalid log.level: '%s' (must be 'DEBUG', 'INFO', 'WARNING' or 'ERROR')", c.Logging.Level)
	}

	logHandlers := []string{"stdout", "file", "both"}
	if !slices.Contains(logHandlers, c.Logging.Handler) {
		return fmt.Errorf("invalid log.handler: '%s' (must be 'stdout', 'file' or 'both')", c.Logging.Handler)
	}

	logFormats := []string{"text", "json"}

	if !slices.Contains(logFormats, c.Logging.StdoutFormat) {
		return fmt.Errorf("invalid log.stdout_foramt: '%s' (must be 'text' or 'json')", c.Logging.StdoutFormat)
	}
	if !slices.Contains(logFormats, c.Logging.FileFormat) {
		return fmt.Errorf("invalid log.file_format: '%s' (must be 'text' or 'json')", c.Logging.FileFormat)
	}

	return nil
}

func (c *ExporterConfiguration) setDefaults() {
	// --- Global ---
	if c.Global.ShutdownTimeout == 0 {
		c.Global.ShutdownTimeout = DefaultGlobalShutdownTimeout
	}
	if c.Global.NumSubordinatesRecords == nil {
		c.Global.NumSubordinatesRecords = []string{}
	}

	// --- HTTP ---
	if c.HTTP.ListenAddress == "" {
		c.HTTP.ListenAddress = DefaultHTTPListenAdderss
	}
	if c.HTTP.MetricsPath == "" {
		c.HTTP.MetricsPath = DefaultHTTPMetricsPath
	}
	if c.HTTP.ReadTimeout == 0 {
		c.HTTP.ReadTimeout = DefaultHTTPReadTimeout
	}
	if c.HTTP.WriteTimeout == 0 {
		c.HTTP.WriteTimeout = DefaultHTTPWriteTimeout
	}
	if c.HTTP.IdleTimeout == 0 {
		c.HTTP.IdleTimeout = DefaultHTTPIdleTimeout
	}
	if c.HTTP.InitialReadTimeout == 0 {
		c.HTTP.InitialReadTimeout = DefaultHTTPInitialReadTimeout
	}

	// --- LDAP ---
	if c.LDAP.PoolConnLimit == 0 {
		c.LDAP.PoolConnLimit = DefaultLDAPPoolConnLimit
	}

	// --- Logging ---
	if c.Logging.Level == "" {
		c.Logging.Level = DefaultLogLevel
	}
	if c.Logging.Handler == "" {
		c.Logging.Handler = DefaultLogHandler
	}
	if c.Logging.File == "" {
		c.Logging.File = DefaultLogFile
	}
	if c.Logging.StdoutFormat == "" {
		c.Logging.StdoutFormat = DefaultLogStdoutFormat
	}
	if c.Logging.FileFormat == "" {
		c.Logging.FileFormat = DefaultLogFileFormat
	}
}

// ReadConfig function reads the configuration from the provided yaml file
// and returns it as a LdapConfiguration structure.
func ReadConfig(configFilePath string) (ExporterConfiguration, error) {
	yamlFile, err := os.ReadFile(configFilePath)
	configuration := ExporterConfiguration{
		Global: globalConfig{
			ShutdownTimeout:        DefaultGlobalShutdownTimeout,
			NumSubordinatesRecords: []string{},
		},
		HTTP: httpConfig{
			ListenAddress:      DefaultHTTPListenAdderss,
			MetricsPath:        DefaultHTTPMetricsPath,
			ReadTimeout:        DefaultHTTPReadTimeout,
			WriteTimeout:       DefaultHTTPWriteTimeout,
			IdleTimeout:        DefaultHTTPIdleTimeout,
			InitialReadTimeout: DefaultHTTPInitialReadTimeout,
		},
		LDAP: ldapConfig{
			PoolConnLimit: DefaultLDAPPoolConnLimit,
		},
		Logging: loggingConfig{
			Level:        DefaultLogLevel,
			Handler:      DefaultLogHandler,
			File:         DefaultLogFile,
			StdoutFormat: DefaultLogStdoutFormat,
			FileFormat:   DefaultLogFileFormat,
		},
	}

	if err != nil {
		return ExporterConfiguration{}, fmt.Errorf("unable to open configuration file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		return ExporterConfiguration{}, fmt.Errorf("yaml unmarshal error: %w", err)
	}
	configuration.setDefaults()
	return configuration, nil
}
