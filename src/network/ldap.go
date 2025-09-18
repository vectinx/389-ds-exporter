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
	config           LdapConnectionPoolConfig
	connectionsCh    chan LdapConn
	totalConnections int
	mu               sync.Mutex
	closing          bool
	doneCh           chan struct{}
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
		config:           cfg,
		connectionsCh:    make(chan LdapConn, cfg.MaxConnections),
		doneCh:           make(chan struct{}),
		totalConnections: 0,
	}

	return pool
}

// Get function gives a connection from the pool. If specified timeout expires, returns an error.
func (pool *LdapConnectionPool) Get(timeout time.Duration) (LdapConn, error) {
	pool.mu.Lock()
	if pool.closing {
		pool.mu.Unlock()

		return nil, errors.New("pool is currently closing")
	}
	pool.mu.Unlock()
	for {
		select {
		case conn := <-pool.connectionsCh:
			if !ldapConnIsAlive(conn) {
				err := conn.Unbind()
				if err != nil {
					slog.Debug("Error closing pooled ldap connection", "err", err)
				}
				pool.decreaseTotalConnections()

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
					pool.decreaseTotalConnections()

					return nil, fmt.Errorf("error creating new connection: %w", err)
				}

				return conn, nil
			}
			pool.mu.Unlock()
			select {
			case conn := <-pool.connectionsCh:
				if !ldapConnIsAlive(conn) {
					err := conn.Unbind()
					if err != nil {
						slog.Debug("Error closing pooled ldap connection", "err", err)
					}
					pool.decreaseTotalConnections()

					continue
				}

				return conn, nil
			case <-time.After(timeout):
				return nil, fmt.Errorf("error getting pooled connection: timeout (%s)", timeout.String())
			}
		}
	}
}

// Put function returns specified connection to pool.
func (pool *LdapConnectionPool) Put(conn LdapConn) {
	if conn == nil {
		return
	}

	select {
	case pool.connectionsCh <- conn:
		pool.mu.Lock()
		if len(pool.connectionsCh) == int(pool.totalConnections) && pool.closing {
			close(pool.doneCh)
		}
		pool.mu.Unlock()
	default:
		err := conn.Unbind()
		if err != nil {
			slog.Debug("Error closing pooled ldap connection", "err", err)
		}
		pool.decreaseTotalConnections()
	}
}

// Close function prevents receiving connections from the pool,
// waits until all connections are returned to the pool and closes them.
func (pool *LdapConnectionPool) Close(ctx context.Context) error {
	pool.mu.Lock()

	pool.closing = true
	canCloseImideatly := len(pool.connectionsCh) == int(pool.totalConnections)
	pool.mu.Unlock()

	if !canCloseImideatly {
		select {
		case <-pool.doneCh:
			// all connections returned to pool
		case <-ctx.Done():
			return fmt.Errorf("timeout while waiting for connections to return: %w", ctx.Err())
		}
	}

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

// decreaseTotalConnections function decreases count of connections, managed by pool.
func (pool *LdapConnectionPool) decreaseTotalConnections() {
	pool.mu.Lock()
	pool.totalConnections--
	pool.mu.Unlock()
}
