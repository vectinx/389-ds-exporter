package network

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// LdapConnectionPoolConfig struct implements LdapConnectionPool configuration.
type LdapConnectionPoolConfig struct {
	ServerURL      string // URL of LDAP Server
	BindDN         string // LDAP server bind DN
	BindPw         string // LDAP server bind Password
	MaxConnections int    // Limit of connections in pool
	DialFunc       func(url string) (LdapConn, error)
}

// LdapConnectionPool implements a pool that manages ldap connections.
type LdapConnectionPool struct {
	config                 LdapConnectionPoolConfig
	connectionsCh          chan LdapConn
	managedConnectionCount int
	mu                     sync.Mutex
	closed                 bool
	doneCh                 chan struct{}
}

// ldapConnIsAlive function checks if specified connection is alive.
func ldapConnIsAlive(conn LdapConn) bool {
	req := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		1, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	_, err := conn.Search(req)

	return err == nil
}

// NewLdapConnectionPool function creates and returns new LdapConnectionPool object with specified config.
func NewLdapConnectionPool(cfg LdapConnectionPoolConfig) *LdapConnectionPool {
	pool := &LdapConnectionPool{
		config:                 cfg,
		connectionsCh:          make(chan LdapConn, cfg.MaxConnections),
		doneCh:                 make(chan struct{}),
		managedConnectionCount: 0,
		closed:                 false,
	}

	return pool
}

// Search function gives a connection from the pool. If specified timeout expires, returns an error.
func (pool *LdapConnectionPool) Search(req *ldap.SearchRequest, timeout time.Duration) (*ldap.SearchResult, error) {

	pool.mu.Lock()
	if pool.closed {
		pool.mu.Unlock()

		return nil, errors.New("pool closed")
	}

	pool.mu.Unlock()

	for {
		select {
		case conn := <-pool.connectionsCh:
			slog.Debug("There is an available connection in the pool, using it")

			if !ldapConnIsAlive(conn) {
				slog.Debug("The connection is broken, deleting it")
				_ = conn.Unbind()
				pool.mu.Lock()
				pool.managedConnectionCount--
				pool.mu.Unlock()

				continue
			}

			return pool.searchAndReturn(conn, req)

		default:
			slog.Debug("There are no connections available in the pool")
			pool.mu.Lock()
			if pool.managedConnectionCount < pool.config.MaxConnections {
				pool.managedConnectionCount++
				pool.mu.Unlock()

				slog.Debug("Creating a new connection in the pool")

				conn, err := pool.newConnection()

				if err != nil {
					pool.mu.Lock()
					pool.managedConnectionCount--
					pool.mu.Unlock()

					return nil, fmt.Errorf("error creating new connection: %w", err)
				}

				return pool.searchAndReturn(conn, req)
			}
			pool.mu.Unlock()

			slog.Debug("Waiting for more connections to become available")
			select {
			case conn := <-pool.connectionsCh:
				slog.Debug("There is an available connection in the pool, using it")
				if !ldapConnIsAlive(conn) {
					_ = conn.Unbind()
					pool.mu.Lock()
					pool.managedConnectionCount--
					pool.mu.Unlock()
					slog.Debug("The connection is broken, deleting it")

					continue
				}

				return pool.searchAndReturn(conn, req)

			case <-time.After(timeout):

				return nil, fmt.Errorf("error getting pooled connection: timeout (%s)", timeout.String())
			}
		}
	}
}

// Close function prevents receiving connections from the pool,
// waits until all connections are returned to the pool and closes them.
func (pool *LdapConnectionPool) Close(ctx context.Context) error {
	pool.mu.Lock()
	if pool.closed {
		pool.mu.Unlock()

		return errors.New("closing of closed pool")
	}
	pool.mu.Unlock()

	pool.mu.Lock()
	pool.closed = true
	pool.mu.Unlock()

	close(pool.connectionsCh)

	hasCloseErrors := false
	errorsCount := 0
	for conn := range pool.connectionsCh {
		if conn.Unbind() != nil {
			hasCloseErrors = true
			errorsCount++
		}
	}

	if hasCloseErrors {
		return fmt.Errorf("pool closed incorrectly - failed to close %v connections", errorsCount)
	}

	return nil
}

// newConnection function creates a new connection to ldap.
func (pool *LdapConnectionPool) newConnection() (LdapConn, error) {

	conn, err := pool.config.DialFunc(pool.config.ServerURL)
	if err == nil {
		err = conn.Bind(pool.config.BindDN, pool.config.BindPw)
		if err == nil {
			return conn, nil
		}
		err := conn.Unbind()
		if err != nil {
			slog.Debug("Error closing pooled ldap connection", "err", err)
		}
	}

	return nil, err
}

// searchAndReturn performs a search through the connection and returns it to the pool
func (pool *LdapConnectionPool) searchAndReturn(conn LdapConn, req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	res, err := conn.Search(req)

	// Try to return connection back to pool
	select {
	case pool.connectionsCh <- conn:
		slog.Debug("The connection was successfully returned to the pool after use")
	default:
		slog.Debug("Error returning connection to pool - removing it")
		_ = conn.Unbind()
		pool.mu.Lock()
		pool.managedConnectionCount--
		pool.mu.Unlock()
	}

	return res, err
}
