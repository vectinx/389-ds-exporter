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
	DefaultGlobalShutdownTimeout int    = 5
	DefaultLDAPTlsSkipVerify     bool   = false
	DefaultLDAPPoolConnLimit     int    = 4
	DefaultLDAPPoolGetTimeout    int    = 5
	DefaultLDAPPoolIdleTime      int    = 300
	DefaultLDAPPoolLifeTime      int    = 3600
	DefaultLDAPDialTimeout       int    = 3
	DefaultCollectorsDefault     string = "standard"

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
	if r.LDAPTlsSkipVerify != nil {
		cfg.LDAPTlsSkipVerify = *r.LDAPTlsSkipVerify
	} else {
		cfg.LDAPTlsSkipVerify = DefaultLDAPTlsSkipVerify
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
