package connections

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-ldap/ldap/v3"
)

var (
	ErrPoolClosed    = errors.New("pool closed")
	ErrBadConnection = errors.New("bad connection")
)

type Conn struct {
	pool *LDAPPool
	conn *pooledConn
}

func (c *Conn) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	res, err := c.conn.conn.Search(req)
	if err != nil && isTransportError(err) {
		c.conn.markBad()
	}
	return res, err
}

func (c *Conn) Close() {
	c.pool.putConn(c.conn)
}

type pooledConn struct {
	pool       *LDAPPool
	sync.Mutex // protects the following fields
	conn       LdapConn
	closed     bool
	bad        atomic.Bool
	inUse      bool
	createdAt  time.Time
	returnedAt time.Time
}

func (pc *pooledConn) markBad() {
	pc.bad.Store(true)
}

func (pc *pooledConn) close() error {
	pc.Lock()
	if pc.closed {
		pc.Unlock()
		return nil
	}
	pc.closed = true
	conn := pc.conn
	pc.Unlock()

	// close the underlying LDAP connection outside of the lock
	_ = conn.Unbind()
	_ = conn.Close()

	pc.pool.mu.Lock()
	pc.pool.numOpen--
	pc.pool.mu.Unlock()

	return nil
}

func (pc *pooledConn) expired(timeout time.Duration) bool {
	if timeout <= 0 {
		return false
	}
	return pc.createdAt.Add(timeout).Before(time.Now())
}

func isTransportError(err error) bool {
	var le *ldap.Error
	if errors.As(err, &le) {
		if le.ResultCode == ldap.ErrorNetwork {
			return true
		}
		if le.Err != nil {
			return true
		}
	}
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}

// LdapConn defines an interface for interacting with an LDAP server.
// It includes basic operations such as binding, searching, and unbinding.
type LdapConn interface {
	Bind(LDAPAuthConfig) error
	Search(*ldap.SearchRequest) (*ldap.SearchResult, error)
	Unbind() error
	Close() error
}

type LDAPAuthConfig struct {
	URL           string
	BindDN        string
	BindPw        string
	TlsSkipVerify bool
	DialTimeout   time.Duration
}

type LDAPPoolConfig struct {
	Auth           LDAPAuthConfig
	MaxConnections int
	MaxIdleTime    time.Duration
	MaxLifeTime    time.Duration
	DialTimeout    time.Duration
	ConnFactory    func(*LDAPAuthConfig) (LdapConn, error)
}

type LDAPPool struct {
	cfg LDAPPoolConfig

	mu           sync.Mutex // protects the following fields
	freeConns    []*pooledConn
	numOpen      int
	maxOpen      int           // <= 0 means unlimited
	maxLifetime  time.Duration // maximum amount of time a connection may be reused
	maxIdleTime  time.Duration // maximum amount of time a connection may be idle before being closed
	closed       bool
	connRequests connRequestSet

	cleanerCh chan struct{}

	waitCount           atomic.Int64 // number of times, the Conn request waited for free connections.
	lifeTimeClosedCount atomic.Int64 // number of connections closed by life time expired
	idleTimeClosedCount atomic.Int64 // number of connections closed by idle time expired
	waitDuration        atomic.Int64

	connFactory func(*LDAPAuthConfig) (LdapConn, error)
}

type LDAPPoolStat struct {
	Open           int
	ClosedIdleTime int
	ClosedLifeTime int
	WaitCount      int
	WaitDuration   int
}

// connReuseStrategy determines how (*pool).conn returns database connections.
type connReuseStrategy uint8

const (
	// alwaysNewConn forces a new connection to the database.
	alwaysNewConn connReuseStrategy = iota
	// cachedOrNewConn returns a cached connection, if available, else waits
	// for one to become available (if MaxOpenConns has been reached) or
	// creates a new database connection.
	cachedOrNewConn
)

// Deprecated: connection opener channel is not used.

func NewLDAPPool(cfg LDAPPoolConfig) *LDAPPool {
	_, cancel := context.WithCancel(context.Background())
	_ = cancel // context cancel is not used currently

	pool := &LDAPPool{
		cfg:         cfg,
		connFactory: cfg.ConnFactory,
		maxLifetime: cfg.MaxLifeTime,
		maxIdleTime: cfg.MaxIdleTime,
		maxOpen:     cfg.MaxConnections,
		cleanerCh:   nil,
	}

	return pool
}

// Conn returns a single connection by either opening a new connection
// or returning an existing connection from the connection pool. Conn will
// block until either a connection is returned or ctx is canceled.
// Queries run on the same Conn will be run in the same database session.
//
// Every Conn must be returned to the database pool after use by
// calling [Conn.Close].
func (pool *LDAPPool) Conn(ctx context.Context) (*Conn, error) {
	var pc *pooledConn
	var err error

	err = pool.retry(func(strategy connReuseStrategy) error {
		pc, err = pool.conn(strategy, ctx)
		return err
	})

	if err != nil {
		return nil, err
	}

	conn := &Conn{
		pool: pool,
		conn: pc,
	}
	pool.mu.Lock()
	slog.Debug("pool counters",
		"max", pool.maxOpen,
		"open", pool.numOpen,
		"waited", pool.waitCount.Load(),
		"lifetime", pool.lifeTimeClosedCount.Load(),
		"idletime", pool.idleTimeClosedCount.Load(),
	)
	pool.mu.Unlock()
	return conn, nil
}

// Close closes the pool, rejecting new acquisitions, waking waiters, and
// closing all idle connections. Active connections will be closed when
// returned via Conn.Close().
func (pool *LDAPPool) Close() error {
	pool.mu.Lock()
	if pool.closed {
		pool.mu.Unlock()
		return nil
	}
	pool.closed = true

	// stop cleaner
	if pool.cleanerCh != nil {
		close(pool.cleanerCh)
		pool.cleanerCh = nil
	}

	// close and clear idle connections
	idle := pool.freeConns
	pool.freeConns = nil

	// wake all waiters with closed signal
	pool.connRequests.CloseAndRemoveAll()

	pool.mu.Unlock()

	for _, pc := range idle {
		if pc != nil {
			_ = pc.close()
		}
	}
	return nil
}

// Stat returns pool usage statistics.
func (pool *LDAPPool) Stat() LDAPPoolStat {
	stat := LDAPPoolStat{}

	pool.mu.Lock()
	stat.Open = pool.numOpen
	pool.mu.Unlock()

	stat.WaitCount = int(pool.waitCount.Load())
	stat.WaitDuration = int(pool.waitDuration.Load())
	stat.ClosedIdleTime = int(pool.idleTimeClosedCount.Load())
	stat.ClosedLifeTime = int(pool.lifeTimeClosedCount.Load())

	return stat
}

//nolint:gocognit,nestif
func (pool *LDAPPool) conn(strategy connReuseStrategy, ctx context.Context) (*pooledConn, error) {
	pool.mu.Lock()
	if pool.closed {
		pool.mu.Unlock()
		return nil, ErrPoolClosed
	}
	pool.mu.Unlock()

	// Check if context expired
	err := ctx.Err()
	if err != nil {
		return nil, err
	}

	pool.mu.Lock()

	// If we have idle connections that we can use
	last := len(pool.freeConns) - 1
	if last >= 0 && strategy == cachedOrNewConn {
		pc := pool.freeConns[last]
		pool.freeConns = pool.freeConns[:last]
		pc.inUse = true

		if pc.bad.Load() {
			pool.mu.Unlock()
			pool.lifeTimeClosedCount.Add(1)
			_ = pc.close()
			return nil, ErrBadConnection
		}

		// lifetime expiration check
		if pc.expired(pool.maxLifetime) {
			pool.mu.Unlock()
			pool.lifeTimeClosedCount.Add(1)
			_ = pc.close()
			return nil, ErrBadConnection
		}

		// idle expiration check
		if pool.maxIdleTime > 0 && pc.returnedAt.Add(pool.maxIdleTime).Before(time.Now()) {
			pool.mu.Unlock()
			pool.idleTimeClosedCount.Add(1)
			_ = pc.close()
			return nil, ErrBadConnection
		}

		pool.mu.Unlock()
		return pc, nil
	}

	if pool.maxOpen > 0 && pool.numOpen >= pool.maxOpen {
		// If there are no idle connections, and we cannott create new one
		req := make(chan connRequest, 1)
		delHandle := pool.connRequests.Add(req)
		pool.waitCount.Add(1)
		pool.mu.Unlock()

		waitStart := time.Now()
		select {
		case <-ctx.Done():
			// Remove the connection request and ensure no value has been sent
			// on it after removing.
			pool.mu.Lock()
			deleted := pool.connRequests.Delete(delHandle)
			pool.mu.Unlock()

			pool.waitDuration.Add(int64(time.Since(waitStart)))

			// If we failed to delete it, that means either the LDAPPool was closed or
			// something else grabbed it and is about to send on it.
			if !deleted {
				select {
				default:
				case ret, ok := <-req:
					if ok && ret.conn != nil {
						pool.putConn(ret.conn)
						// Тут надо вернуть соединение в пул
					}
				}
			}
			return nil, ctx.Err()
		case ret, ok := <-req:
			pool.waitDuration.Add(int64(time.Since(waitStart)))

			if !ok {
				// The req channel can only be closed by method connRequestSet.CloseAndRemoveAll,
				// so if it is closed, we consider the pool closed.
				return nil, ErrPoolClosed
			}

			// Only check if the connection is expired if the strategy is cachedOrNewConns.
			// If we require a new connection, just re-use the connection without looking
			// at the expiry time. If it is expired, it will be checked when it is placed
			// back into the connection pool.
			// This prioritizes giving a valid connection to a client over the exact connection
			// lifetime, which could expire exactly after this point anyway.
			if strategy == cachedOrNewConn && ret.err == nil {
				if (ret.conn != nil && ret.conn.bad.Load()) || ret.conn.expired(pool.maxLifetime) {
					pool.lifeTimeClosedCount.Add(1)
					_ = ret.conn.close()
					return nil, ErrBadConnection
				}
				if pool.maxIdleTime > 0 && ret.conn.returnedAt.Add(pool.maxIdleTime).Before(time.Now()) {
					pool.idleTimeClosedCount.Add(1)
					_ = ret.conn.close()
					return nil, ErrBadConnection
				}
			}
			if ret.conn == nil {
				return nil, ret.err
			}

			return ret.conn, ret.err
		}
	}

	// If we can issue new connection
	pool.numOpen++ // optimistically
	pool.mu.Unlock()

	lc, err := pool.connFactory(&pool.cfg.Auth)
	if err != nil {
		pool.mu.Lock()
		pool.numOpen-- // correct for earlier optimism
		pool.mu.Unlock()
		return nil, fmt.Errorf("dial failed: %w", err)
	}
	err = lc.Bind(pool.cfg.Auth)
	if err != nil {
		_ = lc.Close()
		pool.mu.Lock()
		pool.numOpen-- // correct for earlier optimism
		pool.mu.Unlock()
		return nil, fmt.Errorf("bind failed: %w", err)
	}
	conn := &pooledConn{
		pool:       pool,
		conn:       lc,
		createdAt:  time.Now(),
		returnedAt: time.Now(),
		inUse:      true,
	}
	return conn, nil
}

func (pool *LDAPPool) putConn(pc *pooledConn) {
	var err error
	pool.mu.Lock()
	if !pc.inUse {
		// Idempotent close: ignore duplicate returns
		pool.mu.Unlock()
		return
	}

	if pc.bad.Load() || pc.expired(pool.maxLifetime) {
		pool.lifeTimeClosedCount.Add(1)
		err = ErrBadConnection
	}

	pc.inUse = false
	pc.returnedAt = time.Now()

	if errors.Is(err, ErrBadConnection) {
		// Don't reuse bad connections.
		// Since the conn is considered bad and is being discarded, treat it
		// as closed. Don't decrement the open count here, finalClose will
		// take care of that.
		// pool.maybeOpenNewConnections()
		pool.mu.Unlock()
		_ = pc.close()
		return
	}

	added := pool.putConnLocked(pc, nil)
	pool.mu.Unlock()

	if !added {
		_ = pc.close()
		return
	}
}

// Satisfy a connRequest or put the driverConn in the idle pool and return true
// or return false.
// putConnLocked will satisfy a connRequest if there is one, or it will
// return the *driverConn to the freeConn list if err == nil and the idle
// connection limit will not be exceeded.
// If err != nil, the value of dc is ignored.
// If err == nil, then dc must not equal nil.
// If a connRequest was fulfilled or the *driverConn was placed in the
// freeConn list, then true is returned, otherwise false is returned.
func (pool *LDAPPool) putConnLocked(pc *pooledConn, err error) bool {
	if pool.closed {
		return false
	}
	if pool.maxOpen > 0 && pool.numOpen > pool.maxOpen {
		return false
	}
	if req, ok := pool.connRequests.TakeRandom(); ok {
		if err == nil {
			pc.inUse = true
		}
		req <- connRequest{
			conn: pc,
			err:  err,
		}
		return true
	} else if err == nil && !pool.closed {
		pool.freeConns = append(pool.freeConns, pc)
		pool.startCleanerLocked()
		return true
	}
	return false
}

// maxBadConnRetries is the number of maximum retries if the driver returns
// ErrBadConnection to signal a broken connection before forcing a new
// connection to be opened.
const maxBadConnRetries = 2

func (pool *LDAPPool) retry(fn func(strategy connReuseStrategy) error) error {
	for i := int64(0); i < maxBadConnRetries; i++ {
		err := fn(cachedOrNewConn)
		// retry if err is ErrBadConnection
		if err == nil || !errors.Is(err, ErrBadConnection) {
			return err
		}
	}

	return fn(alwaysNewConn)
}

func (pool *LDAPPool) shortestIdleTimeLocked() time.Duration {
	if pool.maxIdleTime <= 0 {
		return pool.maxLifetime
	}
	if pool.maxLifetime <= 0 {
		return pool.maxIdleTime
	}
	return min(pool.maxIdleTime, pool.maxLifetime)
}

// startCleanerLocked starts connectionCleaner if needed.
func (pool *LDAPPool) startCleanerLocked() {
	if (pool.maxLifetime > 0 || pool.maxIdleTime > 0) && pool.numOpen > 0 && pool.cleanerCh == nil {
		pool.cleanerCh = make(chan struct{}, 1)
		go pool.connectionCleaner(pool.shortestIdleTimeLocked())
	}
}

func (pool *LDAPPool) connectionCleaner(d time.Duration) {
	const minInterval = time.Second

	if d < minInterval {
		d = minInterval
	}
	t := time.NewTimer(d)

	for {
		// capture channel under lock to avoid races with Close() modifying cleanerCh
		pool.mu.Lock()
		ch := pool.cleanerCh
		pool.mu.Unlock()

		select {
		case <-t.C:
		case <-ch: // pool was closed.
		}
		pool.mu.Lock()

		d = pool.shortestIdleTimeLocked()
		if pool.closed || pool.numOpen == 0 || d <= 0 {
			pool.cleanerCh = nil
			pool.mu.Unlock()
			return
		}

		d, closing := pool.connectionCleanerRunLocked(d)
		pool.mu.Unlock()
		for _, c := range closing {
			_ = c.close()
		}

		if d < minInterval {
			d = minInterval
		}

		if !t.Stop() {
			select {
			case <-t.C:
			default:
			}
		}
		t.Reset(d)
	}
}

// connectionCleanerRunLocked removes connections that should be closed from
// freeConn and returns them along side an updated duration to the next check
// if a quicker check is required to ensure connections are checked appropriately.
func (pool *LDAPPool) connectionCleanerRunLocked(d time.Duration) (time.Duration, []*pooledConn) {
	var idleClosing int64
	var closing []*pooledConn
	if pool.maxIdleTime > 0 {
		// As freeConn is ordered by returnedAt process
		// in reverse order to minimise the work needed.
		idleSince := time.Now().Add(-pool.maxIdleTime)
		last := len(pool.freeConns) - 1
		for i := last; i >= 0; i-- {
			c := pool.freeConns[i]
			if c.returnedAt.Before(idleSince) {
				i++
				closing = pool.freeConns[:i:i]
				pool.freeConns = pool.freeConns[i:]
				idleClosing = int64(len(closing))
				pool.idleTimeClosedCount.Add(idleClosing)
				break
			}
		}

		if len(pool.freeConns) > 0 {
			c := pool.freeConns[0]
			if d2 := c.returnedAt.Sub(idleSince); d2 < d {
				// Ensure idle connections are cleaned up as soon as
				// possible.
				d = d2
			}
		}
	}

	if pool.maxLifetime > 0 {
		expiredSince := time.Now().Add(-pool.maxLifetime)
		// Compact in-place preserving order: write index w
		w := 0
		for r := 0; r < len(pool.freeConns); r++ {
			c := pool.freeConns[r]
			if c.createdAt.Before(expiredSince) {
				closing = append(closing, c)
				continue
			}
			// update next deadline
			if d2 := c.createdAt.Sub(expiredSince); d2 < d {
				d = d2
			}
			pool.freeConns[w] = pool.freeConns[r]
			w++
		}
		// zero tail for GC and reslice
		for i := w; i < len(pool.freeConns); i++ {
			pool.freeConns[i] = nil
		}
		pool.freeConns = pool.freeConns[:w]
		pool.lifeTimeClosedCount.Add(int64(len(closing)) - idleClosing)
	}

	return d, closing
}

// connRequest represents one request for a new connection
// When there are no idle connections available, LDAPPool.conn will create
// a new connRequest and put it on the LDAPPool.connRequests list.
type connRequest struct {
	conn *pooledConn
	err  error
}

// connRequestSet is a set of chan connRequest that's
// optimized for:
//
//   - adding an element
//   - removing an element (only by the caller who added it)
//   - taking (get + delete) a random element
type connRequestSet struct {
	// s are the elements in the set.
	s []connRequestAndIndex
}

type connRequestAndIndex struct {
	// req is the element in the set.
	req chan connRequest

	// curIdx points to the current location of this element in
	// connRequestSet.s. It gets set to -1 upon removal.
	curIdx *int
}

// CloseAndRemoveAll closes all channels in the set
// and clears the set.
func (s *connRequestSet) CloseAndRemoveAll() {
	for _, v := range s.s {
		*v.curIdx = -1
		close(v.req)
	}
	s.s = nil
}

// Len returns the length of the set.
func (s *connRequestSet) Len() int { return len(s.s) }

// connRequestDelHandle is an opaque handle to delete an
// item from calling Add.
type connRequestDelHandle struct {
	idx *int // pointer to index; or -1 if not in slice
}

// Add adds v to the set of waiting requests.
// The returned connRequestDelHandle can be used to remove the item from
// the set.
func (s *connRequestSet) Add(v chan connRequest) connRequestDelHandle {
	idx := len(s.s)
	idxPtr := &idx
	s.s = append(s.s, connRequestAndIndex{v, idxPtr})
	return connRequestDelHandle{idxPtr}
}

// Delete removes an element from the set.
//
// It reports whether the element was deleted. (It can return false if a caller
// of TakeRandom took it meanwhile, or upon the second call to Delete).
func (s *connRequestSet) Delete(h connRequestDelHandle) bool {
	idx := *h.idx
	if idx < 0 {
		return false
	}
	s.deleteIndex(idx)
	return true
}

// TakeRandom returns and removes a random element from s
// and reports whether there was one to take. (It returns ok=false
// if the set is empty.)
func (s *connRequestSet) TakeRandom() (chan connRequest, bool) {
	if len(s.s) == 0 {
		return nil, false
	}

	pick, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s.s))))
	e := s.s[pick.Int64()]
	s.deleteIndex(int(pick.Int64()))
	return e.req, true
}

func (s *connRequestSet) deleteIndex(idx int) {
	// Mark item as deleted.
	*(s.s[idx].curIdx) = -1
	// Copy last element, updating its position
	// to its new home.
	if idx < len(s.s)-1 {
		last := s.s[len(s.s)-1]
		*last.curIdx = idx
		s.s[idx] = last
	}
	// Zero out last element (for GC) before shrinking the slice.
	s.s[len(s.s)-1] = connRequestAndIndex{}
	s.s = s.s[:len(s.s)-1]
}
