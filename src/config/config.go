package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type LdapConfiguration struct {
	LdapServerUrl string   `yaml:"ldapServerURL"`
	LdapBindDn    string   `yaml:"ldapBindDn"`
	LdapBindPw    string   `yaml:"ldapBindPw"`
	Backends      []string `yaml:"backends"`
}

// ReadConfig function reads the configuration from the provided yaml file and returns it as a LdapConfiguration structure
func ReadConfig(configFilePath string) (LdapConfiguration, error) {
	yamlFile, err := os.ReadFile(configFilePath)
	configuration := LdapConfiguration{}

	if err != nil {
		return LdapConfiguration{}, fmt.Errorf("unable to open configuration file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		return LdapConfiguration{}, fmt.Errorf("error unmarshaling configuration: %w", err)
	}

	return configuration, nil
}
