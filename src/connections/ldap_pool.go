package connections

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/go-ldap/ldap/v3"
)

var (
	ErrPoolClosed             = errors.New("pool closed")
	ErrPoolFull               = errors.New("pool is full")
	ErrPoolGetTimedOut        = errors.New("timed out while trying to get a connection from the pool")
	ErrPoolClosingWaitTimeout = errors.New("timed out while waiting for all connections to be returned to the pool")
	ErrPoolClosingFailed      = errors.New("the pool was closed with errors")
)

type PooledConn struct {
	conn LdapConn
	pool *LdapConnectionPool
	once sync.Once
}

func (c *PooledConn) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	return c.conn.Search(req)
}

func (c *PooledConn) Close() {
	c.once.Do(func() {
		_ = c.pool.put(c.conn)
		c.pool.wg.Done()
	})
}

type LdapConnectionPoolConfig struct {
	ServerURL      string
	BindDN         string
	BindPw         string
	MaxConnections int
	ConnFactory    func(url string) (LdapConn, error)
}

type LdapConnectionPool struct {
	cfg        LdapConnectionPoolConfig
	connCh     chan LdapConn
	totalConns atomic.Int32
	closed     bool
	mu         sync.Mutex
	wg         sync.WaitGroup
}

func NewLdapConnectionPool(config LdapConnectionPoolConfig) *LdapConnectionPool {
	return &LdapConnectionPool{
		cfg:    config,
		connCh: make(chan LdapConn, config.MaxConnections),
	}
}

func (p *LdapConnectionPool) Get(ctx context.Context) (*PooledConn, error) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}
	p.wg.Add(1)
	p.mu.Unlock()

	for {
		conn := p.tryGetFromChan()
		if conn != nil {
			if isConnAlive(conn) {
				return &PooledConn{conn: conn, pool: p}, nil
			}
			p.totalConns.Add(-1)
			_ = conn.Unbind()
		}

		p.mu.Lock()
		canCreate := int(p.totalConns.Load()) < p.cfg.MaxConnections
		if canCreate {
			p.totalConns.Add(1)
		}
		p.mu.Unlock()

		if canCreate {
			conn, _ := p.cfg.ConnFactory(p.cfg.ServerURL)
			if conn != nil {

				err := conn.Bind(p.cfg.BindDN, p.cfg.BindPw)
				if err == nil {
					return &PooledConn{conn: conn, pool: p}, nil
				}
				_ = conn.Unbind()
			}
			p.totalConns.Add(-1)
		}

		select {
		case <-ctx.Done():
			p.wg.Done()
			return nil, ErrPoolGetTimedOut
		default:
		}
	}
}

func (p *LdapConnectionPool) TotalConnections() int {
	return int(p.totalConns.Load())
}

func (p *LdapConnectionPool) ConnsAtPool() int {
	return len(p.connCh)
}

func (p *LdapConnectionPool) Close(ctx context.Context) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()

		return ErrPoolClosed
	}
	p.closed = true

	p.mu.Unlock()
	doneCh := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		// all connections returned to pool
	case <-ctx.Done():

		return ErrPoolClosingWaitTimeout
	}

	close(p.connCh)

	hasCloseErrors := false
	errorsCount := 0
	for conn := range p.connCh {
		if conn.Unbind() != nil {
			hasCloseErrors = true
			errorsCount++
		}
	}

	if hasCloseErrors {
		return ErrPoolClosingFailed
	}

	return nil
}

func (p *LdapConnectionPool) Closed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.closed
}

func (p *LdapConnectionPool) put(conn LdapConn) error {

	select {
	case p.connCh <- conn:

		return nil
	default:
		_ = conn.Unbind()
		slog.Debug("Unable to put connection to pool")

		return ErrPoolFull
	}
}

func (p *LdapConnectionPool) tryGetFromChan() LdapConn {
	select {
	case conn := <-p.connCh:
		return conn
	default:
		return nil
	}
}

func isConnAlive(conn LdapConn) bool {
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
	if err != nil {
		slog.Debug("Error checking connection", "err", err)
	}
	return err == nil
}
