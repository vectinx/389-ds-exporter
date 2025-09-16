package network

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// LdapConnectionPoolConfig struct implements LdapConnectionPool configuration.
type LdapConnectionPoolConfig struct {
	ServerURL              string        // URL of LDAP Server
	BindDN                 string        // LDAP server bind DN
	BindPw                 string        // LDAP server bind Password
	MaxConnections         int           // Limit of connections in pool
	DialTimeout            time.Duration // Network timeout while creating new connection
	RetryCount             int           // Number of attempts to reconnect to an unavailable connection
	RetryDelay             time.Duration // Delay between reconnect attemnts
	ConnectionAliveTimeout time.Duration // Connection alive check timeout
}

// LdapConnectionPool implements a pool that manages ldap connections.
type LdapConnectionPool struct {
	config           LdapConnectionPoolConfig
	connectionsCh    chan *ldap.Conn
	totalConnections int
	mu               sync.Mutex
	closing          bool
	doneCh           chan struct{}
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

	conn.SetTimeout(timeout)
	_, err := conn.Search(req)
	conn.SetTimeout(ldap.DefaultTimeout)

	return err == nil
}

// NewLdapConnectionPool function creates and returns new LdapConnectionPool object with specified config.
func NewLdapConnectionPool(cfg LdapConnectionPoolConfig) *LdapConnectionPool {
	pool := &LdapConnectionPool{
		config:           cfg,
		connectionsCh:    make(chan *ldap.Conn, cfg.MaxConnections),
		doneCh:           make(chan struct{}),
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
				err := conn.Close()
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

					return nil, err
				}

				return conn, nil
			}
			pool.mu.Unlock()
			select {
			case conn := <-pool.connectionsCh:
				if !ldapConnIsAlive(conn, pool.config.ConnectionAliveTimeout) {
					err := conn.Close()
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
func (pool *LdapConnectionPool) Put(conn *ldap.Conn) {
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
		err := conn.Close()
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
		if conn.Close() != nil {
			hasCloseErrors = true
			errorsCount++
		}
	}

	if hasCloseErrors {
		return fmt.Errorf("pool closed incorrectly - failed to close %v connections", errorsCount)
	}

	return nil
}

// newConnection function creates a new connection to ldap with the specified number of retries.
func (pool *LdapConnectionPool) newConnection() (*ldap.Conn, error) {
	var conn *ldap.Conn

	dialer := &net.Dialer{Timeout: pool.config.DialTimeout}

	err := fmt.Errorf(
		"failed to create connection after %v attempts with %v delay",
		pool.config.RetryCount,
		pool.config.RetryDelay,
	)

	// Since retry is the number of connection attempts after the first failed one,
	// we add this first attempt.
	attempts := pool.config.RetryCount + 1
	for range attempts {
		conn, err = ldap.DialURL(pool.config.ServerURL, ldap.DialWithDialer(dialer))
		if err == nil {
			err = conn.Bind(pool.config.BindDN, pool.config.BindPw)
			if err == nil {
				return conn, nil
			}
			err := conn.Close()
			if err != nil {
				slog.Debug("Error closing pooled ldap connection", "err", err)
			}
		}
		time.Sleep(pool.config.RetryDelay)
	}

	return nil, err
}

// decreaseTotalConnections function decreases count of connections, managed by pool.
func (pool *LdapConnectionPool) decreaseTotalConnections() {
	pool.mu.Lock()
	pool.totalConnections--
	pool.mu.Unlock()
}
