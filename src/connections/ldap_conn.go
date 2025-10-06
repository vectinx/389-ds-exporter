package connections

import (
	"context"
	"net"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// LdapConn defines an interface for interacting with an LDAP server.
// It includes basic operations such as binding, searching, and unbinding.
type LdapConn interface {
	Bind(string, string) error
	Search(*ldap.SearchRequest) (*ldap.SearchResult, error)
	Unbind() error
}

// RealLdapConn is a concrete implementation of the LdapConn interface,
// using a real ldap.Conn from the go-ldap library.
type RealLdapConn struct {
	conn *ldap.Conn
}

// Bind authenticates to the LDAP server using the given bind DN and password.
func (c *RealLdapConn) Bind(bind_dn string, bind_pw string) error {
	return c.conn.Bind(bind_dn, bind_pw)
}

// Search executes the given LDAP search request and returns the result.
func (c *RealLdapConn) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	return c.conn.Search(req)
}

// Unbind closes the LDAP connection.
func (c *RealLdapConn) Unbind() error {
	return c.conn.Unbind()
}

// RealConnectionDialUrl establishes a connection to the LDAP server using the given URL.
// It returns an LdapConn interface backed by a real connection, or an error if the connection fails.
func RealConnectionDialUrl(url string, ctx context.Context, defaultTimeout time.Duration) (LdapConn, error) {
	timeout := defaultTimeout

	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining < timeout {
			timeout = remaining
		}
	}

	dialer := &net.Dialer{Timeout: timeout}

	conn, err := ldap.DialURL(url, ldap.DialWithDialer(dialer))
	if err != nil {
		return nil, err
	}

	return &RealLdapConn{
		conn: conn,
	}, nil
}
