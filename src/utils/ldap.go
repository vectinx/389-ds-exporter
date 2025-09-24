package utils

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"389-ds-exporter/src/config"
	"389-ds-exporter/src/connections"
)

// GetLdapBackendType gets backend parameters from ldap and returns them as a BackendType.
func GetLdapBackendType(conn *connections.PooledConn) (config.BackendType, error) {
	searchAttributesRequest := ldap.NewSearchRequest(
		"cn=config,cn=ldbm database,cn=plugins,cn=config",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		"(objectclass=*)",
		[]string{"nsslapd-backend-implement"},
		nil,
	)

	searchResult, err := conn.Search(searchAttributesRequest)
	if err != nil {
		return "", fmt.Errorf("error determining backend type: %w", err)
	}
	return config.BackendType(searchResult.Entries[0].GetAttributeValue("nsslapd-backend-implement")), nil
}
