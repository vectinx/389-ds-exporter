package config

import (
	"log"
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
func ReadConfig(configFilePath string) LdapConfiguration {
	yamlFile, err := os.ReadFile(configFilePath)
	configuration := LdapConfiguration{}

	if err != nil {
		log.Println("Unable to open configuration file")
	}

	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return configuration
}
