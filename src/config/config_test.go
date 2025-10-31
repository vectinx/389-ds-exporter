package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getConf(t *testing.T, file string) *ExporterConfig {
	config, err := ReadConfig(file)
	require.NoError(t, err, "Reading a valid configuration should not result in an error")
	return config
}

func TestLoadValidConfig(t *testing.T) {
	config := getConf(t, "testdata/valid.yml")

	err := config.Validate()
	require.NoError(t, err, "Validating a valid configuration should not result in an error")

	require.Equal(t, config.CollectorsDefault, "none")
	require.Equal(t, config.CollectorsEnabled, []string{"server", "snmp-server"})
	require.Equal(t, config.DSBackendType, "mdb")
	require.Equal(t, config.DSBackendDBs, []string{"userRoot"})
	require.Equal(t, config.DSNumSubordinateRecords, []string{"cn=users,cn=accounts,dc=example,dc=com"})
	require.Equal(t, config.ShutdownTimeout, 5)
	require.Equal(t, config.LDAPServerURL, "ldap://localhost:389")
	require.Equal(t, config.LDAPBindDN, "cn=directory manager")
	require.Equal(t, config.LDAPBindPw, "12345678")
	require.Equal(t, config.LDAPTlsSkipVerify, false)
	require.Equal(t, config.LDAPPoolConnLimit, 5)
	require.Equal(t, config.LDAPPoolGetTimeout, 5)
	require.Equal(t, config.LDAPDialTimeout, 3)
	require.Equal(t, config.LDAPPoolIdleTime, 600)
	require.Equal(t, config.LDAPPoolLifeTime, 3600)
}

func TestDefaultConfigValues(t *testing.T) {
	config := getConf(t, "testdata/default-values.yml")
	err := config.Validate()
	require.NoError(t, err, "Validating a config with default values should not result in an error")

	require.Equal(t, config.CollectorsDefault, DefaultCollectorsDefault)
	require.Equal(t, config.CollectorsEnabled, []string(nil))
	require.Equal(t, config.DSBackendType, "")
	require.Equal(t, config.DSBackendDBs, []string(nil))
	require.Equal(t, config.DSNumSubordinateRecords, []string(nil))
	require.Equal(t, config.ShutdownTimeout, DefaultShutdownTimeout)
	require.Equal(t, config.LDAPServerURL, "ldap://localhost:389")
	require.Equal(t, config.LDAPBindDN, "cn=directory manager")
	require.Equal(t, config.LDAPBindPw, "12345678")
	require.Equal(t, config.LDAPTlsSkipVerify, DefaultLDAPTlsSkipVerify)
	require.Equal(t, config.LDAPPoolConnLimit, DefaultLDAPPoolConnLimit)
	require.Equal(t, config.LDAPPoolGetTimeout, DefaultLDAPPoolGetTimeout)
	require.Equal(t, config.LDAPDialTimeout, DefaultLDAPDialTimeout)
	require.Equal(t, config.LDAPPoolIdleTime, DefaultLDAPPoolIdleTime)
	require.Equal(t, config.LDAPPoolLifeTime, DefaultLDAPPoolLifeTime)
}

func TestNoRequiredConfigValues(t *testing.T) {
	config := getConf(t, "testdata/no-ldap-url.yml")
	err := config.Validate()
	require.ErrorIs(t, err, ErrNoRequiredValue, "Validation configuration without mandatory should fail")
	require.ErrorContains(t, err, "ldap_server_url")

	config = getConf(t, "testdata/no-ldap-binddn.yml")
	err = config.Validate()
	require.ErrorIs(t, err, ErrNoRequiredValue, "Validation configuration without mandatory should fail")
	require.ErrorContains(t, err, "ldap_bind_dn")

	config = getConf(t, "testdata/no-ldap-bindpw.yml")
	err = config.Validate()
	require.ErrorIs(t, err, ErrNoRequiredValue, "Validation configuration without mandatory should fail")
	require.ErrorContains(t, err, "ldap_bind_pw")
}
