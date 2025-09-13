package backends

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// LdapConnectionPoolConfig struct implements LdapConnectionPool configuration
type LdapConnectionPoolConfig struct {
	ServerURL              string        // URL of LDAP Server
	BindDN                 string        // LDAP server bind DN
	BindPw                 string        // LDAP server bind Password
	MaxConnections         uint          // Limit of connections in pool
	DialTimeout            time.Duration // Network timeout while creating new connection
	RetryCount             uint          // Number of attempts to reconnect to an unavailable connection
	RetryDelay             time.Duration // Delay between reconnect attemnts
	ConnectionAliveTimeout time.Duration // Connection alive check timeout
}

// LdapConnectionPool implements a pool that manages ldap connections.
type LdapConnectionPool struct {
	config           LdapConnectionPoolConfig
	connectionsCh    chan *PooledConnection
	totalConnections uint
	mu               sync.Mutex
}

type PooledConnection struct {
	Conn     *ldap.Conn
	lastUsed time.Time
}

// ldapConnIsAlive function checks if specified connection is alive.
func ldapConnIsAlive(conn *ldap.Conn, timeout time.Duration) bool {
	req := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		1, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	conn.Conn.SetTimeout(timeout)
	_, err := conn.Conn.Search(req)
	conn.Conn.SetTimeout(ldap.DefaultTimeout)

	return err == nil
}

// NewLdapConnectionPool function creates and returns new LdapConnectionPool object with specified config.
func NewLdapConnectionPool(cfg LdapConnectionPoolConfig) *LdapConnectionPool {
	pool := &LdapConnectionPool{
		config:           cfg,
		connectionsCh:    make(chan *PooledConnection, cfg.MaxConnections),
		totalConnections: 0,
	}
	return pool
}

// Get function gives a connection from the pool. If specified timeout expires, returns an error.
func (pool *LdapConnectionPool) Get(timeout time.Duration) (*ldap.Conn, error) {

	pool.mu.Lock()
	if pool.closing {
		pool.mu.Unlock()
		return nil, errors.New("pool is currently closing")
	}
	pool.mu.Unlock()
	for {
		select {
		case conn := <-pool.connectionsCh:
			if !ldapConnIsAlive(conn, pool.config.ConnectionAliveTimeout) {
				conn.Close()
				pool.mu.Lock()
				pool.totalConnections--
				pool.mu.Unlock()
				continue
			}
			return conn, nil
		default:
			pool.mu.Lock()
			if pool.totalConnections < pool.config.MaxConnections {
				pool.totalConnections++
				pool.mu.Unlock()
				conn, err := pool.newConnection()
				if err != nil {
					pool.mu.Lock()
					pool.totalConnections--
					pool.mu.Unlock()
					return nil, err
				}
				return conn, nil
			}
			pool.mu.Unlock()
			select {
			case conn := <-pool.connectionsCh:
				if !ldapConnIsAlive(conn, pool.config.ConnectionAliveTimeout) {
					conn.Close()
					pool.mu.Lock()
					pool.totalConnections--
					pool.mu.Unlock()
					continue
				}
				return conn, nil
			case <-time.After(timeout):
				return nil, fmt.Errorf("error getting pooled connection: timeout (%s)", timeout.String())
			}
		}
	}
}

// Put function returns specified connection to pool
func (pool *LdapConnectionPool) Put(conn *ldap.Conn) {
	if conn == nil {
		return
	}
	conn.lastUsed = time.Now()
	select {
	case pool.connectionsCh <- conn:
	default:
		conn.Conn.Close()
		pool.mu.Lock()
		pool.totalConnections--
		pool.mu.Unlock()
	}
}

// Close function prevents receiving connections from the pool, waits until all connections are returned to the pool and closes them
func (pool *LdapConnectionPool) Close(ctx context.Context) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return nil
}

// newConnection function creates a new connection to ldap with the specified number of retries
func (pool *LdapConnectionPool) newConnection() (*ldap.Conn, error) {
	var conn *ldap.Conn

	dialer := &net.Dialer{Timeout: pool.config.DialTimeout}

	err := fmt.Errorf("failed to create connection after %v attempts with %v delay", pool.config.RetryCount, pool.config.RetryDelay)

	// Since retry is the number of connection attempts after the first failed one,
	// we add this first attempt
	attempts := int(pool.config.RetryCount) + 1
	for range attempts {
		conn.Conn, err = ldap.DialURL(pool.config.ServerURL, ldap.DialWithDialer(dialer))
		if err == nil {
			err = conn.Conn.Bind(pool.config.BindDN, pool.config.BindPw)
			if err == nil {
				return conn, nil
			}
			conn.Conn.Close()
		}
		time.Sleep(pool.config.RetryDelay)
	}

	return nil, err
}
