package ldap

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"github.com/go-ldap/ldap/v3"
)

// RealLdapConn is a concrete implementation of the LdapConn interface,
// using a real ldap.Conn from the go-ldap library.
type RealLdapConn struct {
	conn *ldap.Conn
}

// Bind authenticates to the LDAP server using the given bind DN and password.
func (c *RealLdapConn) Bind(auth AuthConfig) error {
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
func RealConnectionDialUrl(auth *AuthConfig) (Conn, error) {
	dialer := &net.Dialer{Timeout: auth.DialTimeout}

	var dialOpts []ldap.DialOpt

	parsed, err := url.Parse(auth.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ldap url: %w", err)
	}

	// We specifically disable the warning "G402 (CWE-295): TLS InsecureSkipVerify may be true",
	// because we specifically leave the option to disable TLS verification.
	if parsed.Scheme == "ldaps" { // #nosec G402
		tlsConfig := &tls.Config{
			InsecureSkipVerify: auth.TlsSkipVerify,
			MinVersion:         tls.VersionTLS12,
		}
		dialOpts = append(dialOpts, ldap.DialWithTLSConfig(tlsConfig))
	}

	dialOpts = append(dialOpts, ldap.DialWithDialer(dialer))

	conn, err := ldap.DialURL(auth.URL, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &RealLdapConn{
		conn: conn,
	}, nil
}
