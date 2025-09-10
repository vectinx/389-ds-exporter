/*
Package controllers provides 'controllers' - structures for obtaining data for metrics
*/
package controllers

import (
	"fmt"
	"log"
	"slices"

	"github.com/go-ldap/ldap/v3"
)

/*
LdapEntryController is the structure that performs the functions of connecting to ldap
and receiving the specified attributes from there.
*/
type LdapEntryController struct {
	ldapServerURL string
	bindDn        string
	bindPassword  string
	baseDn        string
	attributes    []string
}

// NewLdapEntryController function creates new instance of LdapEntryController based on providerd arguments
func NewLdapEntryController(
	ldapServerURL string,
	bindDn string,
	bindPassword string,
	baseDn string,
	attributes []string,
) *LdapEntryController {
	return &LdapEntryController{
		ldapServerURL: ldapServerURL,
		bindDn:        bindDn,
		bindPassword:  bindPassword,
		baseDn:        baseDn,
		attributes:    attributes,
	}
}

// Get the specified attributes from there and returns them as map[string][]string.
func (c *LdapEntryController) Get() (map[string][]string, error) {
	ldapConnection, err := ldap.DialURL(c.ldapServerURL)
	if err != nil {
		return nil, fmt.Errorf("creating ldap connection failed: %w", err)
	}

	defer func() {
		if ldapConnection.Close() != nil {
			log.Printf("Error closing LDAP connection: %s", err)
		}
	}()

	err = ldapConnection.Bind(c.bindDn, c.bindPassword)
	if err != nil {
		return nil, fmt.Errorf("LDAP bind request failed with error: %w", err)
	}

	searchAttributesRequest := ldap.NewSearchRequest(
		c.baseDn,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		"(objectclass=*)",
		c.attributes,
		nil,
	)

	searchResult, err := ldapConnection.Search(searchAttributesRequest)
	if err != nil {
		return nil, fmt.Errorf("LDAP Search request failed with error: %w", err)
	}

	returnValue := make(map[string][]string)

	for _, attr := range searchResult.Entries[0].Attributes {
		if !slices.Contains(c.attributes, attr.Name) {
			continue
		}

		returnValue[attr.Name] = attr.Values
	}

	return returnValue, nil
}
