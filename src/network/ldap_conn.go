package network

import (
	"github.com/go-ldap/ldap/v3"
)

type LdapConn interface {
	Bind(string, string) error
	Search(*ldap.SearchRequest) (*ldap.SearchResult, error)
	Unbind() error
}

type RealLdapConn struct {
	conn *ldap.Conn
}

func (c *RealLdapConn) Bind(bind_dn string, bind_pw string) error {
	return c.conn.Bind(bind_dn, bind_pw)
}

func (c *RealLdapConn) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	return c.conn.Search(req)
}

func (c *RealLdapConn) Unbind() error {
	return c.conn.Unbind()
}

func RealConnectionDialUrl(url string) (LdapConn, error) {
	conn, err := ldap.DialURL(url)
	if err != nil {
		return nil, err
	}

	return &RealLdapConn{
		conn: conn,
	}, nil
}
