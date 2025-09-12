package backends

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
)

type LdapConnectionPoolConfig struct {
	ServerURL              string        // URL of LDAP Server
	BindDN                 string        // LDAP server bind DN
	BindPassword           string        // LDAP server bind Password
	ConnectionsLimit       int           // Limit of connections in pool
	MaxIdleTime            time.Duration // The time after which an unused connection will be considered idle and recreated
	MaxLifeTime            time.Duration // The time after which a connection will be considered ond of life and recreated
	DialTimeout            time.Duration // Network timeout while creating new connection
	RetryCount             int           // Number of attempts to reconnect to an unavailable connection
	RetryDelay             time.Duration // Delay between reconnect attemnts
	ConnectionAliveTimeout time.Duration // Connection alive check timeout
}

// LdapConnectionPool structure providers LDAP connection pool functionality
type LdapConnectionPool struct {
	config            LdapConnectionPoolConfig
	connectionsPool   chan *PooledConnection
	mutex             sync.Mutex
	activeConnections int
}

// PooledConnection struct
type PooledConnection struct {
	Conn     *ldap.Conn
	created  time.Time
	lastUsed time.Time
}

func NewLdapConnectionPool(cfg LdapConnectionPoolConfig) *LdapConnectionPool {
	pool := &LdapConnectionPool{
		config:          cfg,
		connectionsPool: make(chan *PooledConnection, cfg.ConnectionsLimit),
	}
	return pool
}

// Get connection from pool with timeout
func (pool *LdapConnectionPool) Get(timeout time.Duration) (*PooledConnection, error) {

	select {
	// If there are connections in the pool, we try to return the connection.
	case pooledConnection := <-pool.connectionsPool:

		// If the connection MaxLifeTime has expired
		if time.Since(pooledConnection.created) > pool.config.MaxLifeTime {
			log.Print("Expired connection detected. Recreating it ...")
			pooledConnection.Conn.Close()
			pool.decreasActiveConnections()
			return pool.newConnection()
		}
		// If the connection ttl has expired, close it and create a new one.
		if time.Since(pooledConnection.lastUsed) > pool.config.MaxIdleTime {
			log.Print("Idle connection detected. Recreating it ...")
			pooledConnection.Conn.Close()
			pool.decreasActiveConnections()
			return pool.newConnection()
		}

		// If connection is alive - return it.
		if isConnectionAlive(pooledConnection.Conn, pool.config.ConnectionAliveTimeout) {
			log.Print("Reusing already opened connection ...")
			return pooledConnection, nil
		}
		// Else close broken connection and return new.
		log.Print("Broken connection detected. Recreating it ...")
		pooledConnection.Conn.Close()
		pool.decreasActiveConnections()
		return pool.newConnection()
	default:
		log.Print("Pool doesnt contain free connections ...")
		pool.mutex.Lock()
		if pool.activeConnections < pool.config.ConnectionsLimit {
			pool.activeConnections++
			pool.mutex.Unlock()
			log.Print("Creating new connection...")
			return pool.newConnection()
		}
		pool.mutex.Unlock()
		log.Print("Wait for free connection ...")
		select {
		case pooledConnection := <-pool.connectionsPool:
			if isConnectionAlive(pooledConnection.Conn, pool.config.ConnectionAliveTimeout) {
				log.Print("Waited for free connection. Returning it ...")
				return pooledConnection, nil
			}
			pooledConnection.Conn.Close()
			pool.decreasActiveConnections()
			return pool.newConnection()

		case <-time.After(timeout):
			return nil, errors.New("timeout exceeded waiting for connection")
		}
	}
}

// Put connection back to pool
func (pool *LdapConnectionPool) Put(conn *PooledConnection) {
	if conn == nil {
		return
	}
	conn.lastUsed = time.Now()
	select {
	case pool.connectionsPool <- conn:
	default:
		// If the pool is full, close the connection.
		_ = conn.Conn.Close()
		pool.decreasActiveConnections()
	}
}

// Close LdapConnectionsPool
func (pool *LdapConnectionPool) Close() error {
	close(pool.connectionsPool)
	for conn := range pool.connectionsPool {
		conn.Conn.Close()
	}
	return nil
}

func (pool *LdapConnectionPool) decreasActiveConnections() {
	pool.mutex.Lock()
	pool.activeConnections--
	pool.mutex.Unlock()
}

func (pool *LdapConnectionPool) newConnection() (*PooledConnection, error) {
	pooledConnection := &PooledConnection{}
	pooledConnection.created = time.Now()
	var err error
	dialer := &net.Dialer{
		Timeout: pool.config.DialTimeout,
	}
	for i := 0; i < pool.config.RetryCount; i++ {
		pooledConnection.Conn, err = ldap.DialURL(pool.config.ServerURL, ldap.DialWithDialer(dialer))
		if err == nil {
			err = pooledConnection.Conn.Bind(pool.config.BindDN, pool.config.BindPassword)
			if err == nil {
				return pooledConnection, nil
			}
			pooledConnection.Conn.Close()
		}
		time.Sleep(pool.config.RetryDelay)
	}
	pool.decreasActiveConnections()
	return nil, err
}

// isConnectionAlive finction checks if specified connection is opened
func isConnectionAlive(conn *ldap.Conn, timeout time.Duration) bool {
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
