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

	if c.Global.BackendImplement == "" {
		return errors.New("configuration parameter global.ds_backend_implement is required")
	}

	switch c.Global.BackendImplement {
	case BackendBDB, BackendMDB:
		// valid value - pass
	default:
		return fmt.Errorf("invalid global.ds_backend_implement: '%q' (must be 'bdb' or 'mdb')", c.Global.BackendImplement)
	}

	if c.LDAP.GetPoolConnLimit() <= 0 {
		return errors.New("invalid ldap.pool_conn_limit: must be greater than 0")
	}

	logLevels := []string{"DEBUG", "INFO", "WARNING", "ERROR"}
	if !slices.Contains(logLevels, c.Logging.GetLevel()) {
		return fmt.Errorf("invalid log.level: '%s' (must be 'DEBUG', 'INFO', 'WARNING' or 'ERROR')", c.Logging.GetLevel())
	}

	logHandlers := []string{"stdout", "file", "both"}
	if !slices.Contains(logHandlers, c.Logging.GetHandler()) {
		return fmt.Errorf("invalid log.handler: '%s' (must be 'stdout', 'file' or 'both')", c.Logging.GetHandler())
	}

	logFormats := []string{"text", "json"}

	if !slices.Contains(logFormats, c.Logging.GetStdoutFormat()) {
		return fmt.Errorf("invalid log.stdout_foramt: '%s' (must be 'text' or 'json')", c.Logging.GetStdoutFormat())
	}
	if !slices.Contains(logFormats, c.Logging.GetFileFormat()) {
		return fmt.Errorf("invalid log.file_format: '%s' (must be 'text' or 'json')", c.Logging.GetFileFormat())
	}

	return nil
}

// ReadConfig function reads the configuration from the provided yaml file
// and returns it as a LdapConfiguration structure.
func ReadConfig(configFilePath string) (ExporterConfiguration, error) {
	yamlFile, err := os.ReadFile(configFilePath)
	configuration := ExporterConfiguration{}

	if err != nil {
		return ExporterConfiguration{}, fmt.Errorf("unable to open configuration file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		return ExporterConfiguration{}, fmt.Errorf("yaml unmarshal error: %w", err)
	}

	return configuration, nil
}
