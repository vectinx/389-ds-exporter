package connections

import (
	"net"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// RealLdapConn is a concrete implementation of the LdapConn interface,
// using a real ldap.Conn from the go-ldap library.
type RealLdapConn struct {
	conn *ldap.Conn
}

// Bind authenticates to the LDAP server using the given bind DN and password.
func (c *RealLdapConn) Bind(auth LDAPAuthConfig) error {
	return c.conn.Bind(auth.BindDN, auth.BindPw)
}

// Search executes the given LDAP search request and returns the result.
func (c *RealLdapConn) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	return c.conn.Search(req)
}

// Unbind closes the LDAP connection.
func (c *RealLdapConn) Unbind() error {
	return c.conn.Unbind()
}

// Close closes the LDAP connection.
func (c *RealLdapConn) Close() error {
	return c.conn.Close()
}

// RealConnectionDialUrl establishes a connection to the LDAP server using the given URL.
// It returns an LdapConn interface backed by a real connection, or an error if the connection fails.
func RealConnectionDialUrl(auth *LDAPAuthConfig, timeout time.Duration) (LdapConn, error) {
	dialer := &net.Dialer{Timeout: timeout}

	conn, err := ldap.DialURL(auth.URL, ldap.DialWithDialer(dialer))
	if err != nil {
		return nil, err
	}

	return &RealLdapConn{
		conn: conn,
	}, nil
}
