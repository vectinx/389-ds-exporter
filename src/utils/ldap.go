package utils

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-ldap/ldap/v3"

	connections "389-ds-exporter/src/ldap"
)

// GetLdapBackendType gets backend parameters from ldap and returns them as a BackendType.
func GetLdapBackendType(conn *connections.PoolConn) (*string, error) {
	if conn == nil {
		return nil, errors.New("connection is nil")
	}
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
		return nil, fmt.Errorf("error determining backend type: %w", err)
	}

	result := searchResult.Entries[0].GetAttributeValue("nsslapd-backend-implement")
	return &result, nil
}

// GetLdapBackendInstances gets backend instances from ldap and returns them as []string.
func GetLdapBackendInstances(conn *connections.PoolConn) ([]string, error) {
	if conn == nil {
		return nil, errors.New("connection is nil")
	}
	searchAttributesRequest := ldap.NewSearchRequest(
		"cn=ldbm database,cn=plugins,cn=config",
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		"(objectClass=nsBackendInstance)",
		[]string{"cn"},
		nil,
	)

	searchResult, err := conn.Search(searchAttributesRequest)
	if err != nil {
		return nil, fmt.Errorf("error searching for backend instances: %w", err)
	}

	results := []string{}

	for _, entry := range searchResult.Entries {
		cn := entry.GetAttributeValue("cn")
		if cn != "" {
			results = append(results, cn)
		} else {
			slog.Warn("Error getting backend name from record", "err", err, "entry", entry.DN)
		}
	}

	return results, nil
}
