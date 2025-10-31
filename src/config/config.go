package config

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"gopkg.in/yaml.v2"
)

type EnabledCollectorsType string

var (
	ErrNoRequiredValue = errors.New("required field is not specified")
)

const (
	DefaultShutdownTimeout    int    = 5
	DefaultLDAPTlsSkipVerify  bool   = false
	DefaultLDAPPoolConnLimit  int    = 4
	DefaultLDAPPoolGetTimeout int    = 5
	DefaultLDAPPoolIdleTime   int    = 300
	DefaultLDAPPoolLifeTime   int    = 3600
	DefaultLDAPDialTimeout    int    = 3
	DefaultCollectorsDefault  string = "standard"

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

	LDAPServerURL      string `yaml:"ldap_server_url"`
	LDAPBindDN         string `yaml:"ldap_bind_dn"`
	LDAPBindPw         string `yaml:"ldap_bind_pw"`
	LDAPTlsSkipVerify  bool   `yaml:"ldap_tls_skip_verify"`
	LDAPPoolConnLimit  int    `yaml:"ldap_pool_conn_limit"`
	LDAPPoolGetTimeout int    `yaml:"ldap_pool_get_timeout"`
	LDAPPoolIdleTime   int    `yaml:"ldap_pool_idle_time"`
	LDAPPoolLifeTime   int    `yaml:"ldap_pool_life_time"`
	LDAPDialTimeout    int    `yaml:"ldap_dial_timeout"`
}

type rawConfig struct {
	ShutdownTimeout         *int     `yaml:"shutdown_timeout"`
	CollectorsDefault       *string  `yaml:"collectors_default"`
	CollectorsEnabled       []string `yaml:"collectors_enabled"`
	DSNumSubordinateRecords []string `yaml:"ds_numsubordinate_records"`
	DSBackendType           *string  `yaml:"ds_backend_type"`
	DSBackendDBs            []string `yaml:"ds_backend_dbs"`

	LDAPServerURL      *string `yaml:"ldap_server_url"`
	LDAPBindDN         *string `yaml:"ldap_bind_dn"`
	LDAPBindPw         *string `yaml:"ldap_bind_pw"`
	LDAPTlsSkipVerify  *bool   `yaml:"ldap_tls_skip_verify"`
	LDAPPoolConnLimit  *int    `yaml:"ldap_pool_conn_limit"`
	LDAPPoolGetTimeout *int    `yaml:"ldap_pool_get_timeout"`
	LDAPPoolIdleTime   *int    `yaml:"ldap_pool_idle_time"`
	LDAPPoolLifeTime   *int    `yaml:"ldap_pool_life_time"`
	LDAPDialTimeout    *int    `yaml:"ldap_dial_timeout"`
}

func setDefaultIfNotDefined[T any](pointer *T, value *T, defaultValue T) {
	if pointer != nil {
		*value = *pointer
	} else {
		*value = defaultValue
	}
}

func (r *rawConfig) toConfig() *ExporterConfig {
	cfg := &ExporterConfig{}

	// Global
	setDefaultIfNotDefined(r.ShutdownTimeout, &cfg.ShutdownTimeout, DefaultShutdownTimeout)
	setDefaultIfNotDefined(r.CollectorsDefault, &cfg.CollectorsDefault, DefaultCollectorsDefault)

	setDefaultIfNotDefined(r.DSBackendType, &cfg.DSBackendType, "")

	cfg.CollectorsEnabled = r.CollectorsEnabled
	cfg.DSNumSubordinateRecords = r.DSNumSubordinateRecords
	cfg.DSBackendDBs = r.DSBackendDBs

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

	setDefaultIfNotDefined(r.LDAPTlsSkipVerify, &cfg.LDAPTlsSkipVerify, DefaultLDAPTlsSkipVerify)
	setDefaultIfNotDefined(r.LDAPPoolConnLimit, &cfg.LDAPPoolConnLimit, DefaultLDAPPoolConnLimit)
	setDefaultIfNotDefined(r.LDAPPoolGetTimeout, &cfg.LDAPPoolGetTimeout, DefaultLDAPPoolGetTimeout)
	setDefaultIfNotDefined(r.LDAPPoolIdleTime, &cfg.LDAPPoolIdleTime, DefaultLDAPPoolIdleTime)
	setDefaultIfNotDefined(r.LDAPPoolLifeTime, &cfg.LDAPPoolLifeTime, DefaultLDAPPoolLifeTime)
	setDefaultIfNotDefined(r.LDAPDialTimeout, &cfg.LDAPDialTimeout, DefaultLDAPDialTimeout)

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

	if c.LDAPServerURL == "" {
		return fmt.Errorf("ldap_server_url: %w", ErrNoRequiredValue)
	}

	if c.LDAPBindDN == "" {
		return fmt.Errorf("ldap_bind_dn: %w", ErrNoRequiredValue)
	}

	if c.LDAPBindPw == "" {
		return fmt.Errorf("ldap_bind_pw: %w", ErrNoRequiredValue)
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
