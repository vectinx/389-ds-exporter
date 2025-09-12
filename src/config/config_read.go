package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Validate function cheks if provided configuration is valid
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
		return errors.New("configuration parameter global.backend_type is required")
	}

	switch c.Global.BackendImplement {
	case BackendBDB, BackendMDB:
		// valid value - pass
	default:
		return fmt.Errorf("invalid global.backend_type: %q (must be 'bdb' or 'mdb')", c.Global.BackendImplement)
	}

	if c.LDAP.ConnectionPool.GetConnectionsLimit() <= 0 {
		return errors.New("invalid ldap.connection_pool.connections_limit: must be greater than 0")
	}

	return nil
}

// ReadConfig function reads the configuration from the provided yaml file and returns it as a LdapConfiguration structure
func ReadConfig(configFilePath string) (ExporterConfiguration, error) {
	yamlFile, err := os.ReadFile(configFilePath)
	configuration := ExporterConfiguration{}

	if err != nil {
		return ExporterConfiguration{}, fmt.Errorf("unable to open configuration file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		return ExporterConfiguration{}, fmt.Errorf("error unmarshaling configuration: %w", err)
	}

	return configuration, nil
}
