package config

import "gopkg.in/yaml.v2"

type BackendType string
type EnabledCollectorsType string

const (
	BackendBDB         BackendType           = "bdb"
	BackendMDB         BackendType           = "mdb"
	CollectorsAll      EnabledCollectorsType = "all"
	CollectorsStandard EnabledCollectorsType = "standard"
	CollectorsDisabled EnabledCollectorsType = "disabled"
)

// ExporterConfiguration represents exporter configuration top level struct.
type ExporterConfiguration struct {
	Global  globalConfig  `yaml:"global"`
	HTTP    httpConfig    `yaml:"http"`
	LDAP    ldapConfig    `yaml:"ldap"`
	Logging loggingConfig `yaml:"log"`
}

type globalConfig struct {
	ShutdownTimeout        int                   `yaml:"shutdown_timeout"`
	CollectorsDefault      EnabledCollectorsType `yaml:"collectors_default"`
	CollectorsEnabled      []string              `yaml:"collectors_enabled"`
	CollectorsDisabled     []string              `yaml:"collectors_disabled"`
	NumSubordinatesRecords []string              `yaml:"ds_numsubordinate_records"`
}

type httpConfig struct {
	ListenAddress      string `yaml:"listen_address"`
	MetricsPath        string `yaml:"metrics_path"`
	ReadTimeout        int    `yaml:"read_timeout"`
	WriteTimeout       int    `yaml:"write_timeout"`
	IdleTimeout        int    `yaml:"idle_timeout"`
	InitialReadTimeout int    `yaml:"initial_read_timeout"`
}

type ldapConfig struct {
	ServerURL     string `yaml:"server_url"`
	BindDN        string `yaml:"bind_dn"`
	BindPw        string `yaml:"bind_pw"`
	PoolConnLimit int    `yaml:"pool_conn_limit"`
}

type loggingConfig struct {
	Level        string `yaml:"level"`
	Handler      string `yaml:"handler"`
	File         string `yaml:"file"`
	StdoutFormat string `yaml:"stdout_foramt"`
	FileFormat   string `yaml:"file_format"`
}

// String returns string describing config with default values.
func (c *ExporterConfiguration) String() string {
	safeCfg := c

	safeCfg.LDAP.BindPw = "*****"

	out, err := yaml.Marshal(&safeCfg)
	if err != nil {
		return ""
	}
	return string(out)
}
