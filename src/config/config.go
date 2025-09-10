package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ExporterConfiguration struct represents exporter configuration from YAML file
type ExporterConfiguration struct {
	LdapServerUrl string      `yaml:"ldapServerURL"`
	LdapBindDn    string      `yaml:"ldapBindDn"`
	LdapBindPw    string      `yaml:"ldapBindPw"`
	Backends      []string    `yaml:"backends"`
}

// Validate function cheks if provided configuration is valid
func (c *ExporterConfiguration) Validate() error {

	if c.LdapServerUrl == "" {
		return errors.New("configuration parameter LdapServerUrl is required")
	}

	if c.LdapBindDn == "" {
		return errors.New("configuration parameter LdapBindDn is required")
	}

	if c.LdapBindPw == "" {
		return errors.New("configuration parameter LdapBindPw is required")
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
