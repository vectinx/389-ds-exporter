package connections

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// fake LDAP connection for testing.
type fakeLDAP struct {
	bindOK    bool
	searchErr atomic.Value // error
	hasErr    atomic.Bool
	closed    atomic.Bool
}

func (f *fakeLDAP) Bind(_ LDAPAuthConfig) error { f.bindOK = true; return nil }
func (f *fakeLDAP) Search(_ *ldap.SearchRequest) (*ldap.SearchResult, error) {
	if f.hasErr.Load() {
		err, _ := f.searchErr.Load().(error)
		if err != nil {
			return nil, err
		}
	}
	return &ldap.SearchResult{}, nil
}
func (f *fakeLDAP) Unbind() error { f.closed.Store(true); return nil }
func (f *fakeLDAP) Close() error  { f.closed.Store(true); return nil }

func makePool(t *testing.T) *LDAPPool {
	t.Helper()
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) {
		return &fakeLDAP{}, nil
	}
	return NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 1})
}

func TestBadConnMarkingOnTransportError(t *testing.T) {
	pool := makePool(t)
	ctx := context.Background()

	c, err := pool.Conn(ctx)
	if err != nil {
		t.Fatalf("conn: %v", err)
	}

	// inject transport-like error
	le := &ldap.Error{ResultCode: ldap.ErrorNetwork, Err: errors.New("net")}
	c.conn.conn.(*fakeLDAP).searchErr.Store(le)
	c.conn.conn.(*fakeLDAP).hasErr.Store(true)

	_, _ = c.Search(&ldap.SearchRequest{})
	c.Close()

	// next acquire should not reuse bad connection
	c2, err := pool.Conn(ctx)
	if err != nil {
		t.Fatalf("conn2: %v", err)
	}
	if c2.conn == c.conn {
		t.Fatalf("bad connection was reused")
	}
	c2.Close()
}

func TestIdempotentClose(t *testing.T) {
	pool := makePool(t)
	ctx := context.Background()

	c, err := pool.Conn(ctx)
	if err != nil {
		t.Fatalf("conn: %v", err)
	}
	c.Close()
	// second close should not panic
	c.Close()
}

func TestIdleAndLifetimeExpiration(t *testing.T) {
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) {
		return &fakeLDAP{}, nil
	}
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 1})
	// set small times to trigger expiration quickly
	p.maxIdleTime = 10 * time.Millisecond
	p.maxLifetime = 10 * time.Millisecond

	ctx := context.Background()
	c, err := p.Conn(ctx)
	if err != nil {
		t.Fatalf("conn: %v", err)
	}
	c.Close()
	time.Sleep(20 * time.Millisecond)

	// next acquire should create new connection (old expired)
	c2, err := p.Conn(ctx)
	if err != nil {
		t.Fatalf("conn2: %v", err)
	}
	if c2.conn == c.conn {
		t.Fatalf("expired connection was reused")
	}
	c2.Close()
}

func TestPoolClose(t *testing.T) {
	p := makePool(t)
	ctx := context.Background()
	c, err := p.Conn(ctx)

	if err != nil {
		t.Fatalf("conn: %v", err)
	}
	c.Close()

	err = p.Close()
	if err != nil {
		t.Fatalf("close: %v", err)
	}

	_, err = p.Conn(ctx)
	if err == nil {
		t.Fatalf("expected error after pool close")
	}
}

func TestConnContextCanceledImmediately(t *testing.T) {
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) { return &fakeLDAP{}, nil }
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 1})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := p.Conn(ctx)
	if err == nil {
		t.Fatalf("expected context error")
	}
}

func TestWaiterContextTimeout(t *testing.T) {
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) { return &fakeLDAP{}, nil }
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 1})
	baseCtx := context.Background()
	// occupy the only connection
	c, _ := p.Conn(baseCtx)

	// waiter with short timeout
	ctx, cancel := context.WithTimeout(baseCtx, 10*time.Millisecond)
	defer cancel()
	_, err := p.Conn(ctx)
	if err == nil {
		t.Fatalf("expected timeout waiting for connection")
	}
	c.Close()
}

func TestMaxOpenFairnessBasic(t *testing.T) {
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) { return &fakeLDAP{}, nil }
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 2})
	ctx := context.Background()

	// take all
	c1, _ := p.Conn(ctx)
	c2, _ := p.Conn(ctx)

	// spawn 5 waiters
	const waiters = 5
	done := make(chan struct{})
	for i := 0; i < waiters; i++ {
		go func() {
			c, err := p.Conn(ctx)
			if err == nil {
				c.Close()
			}
			done <- struct{}{}
		}()
	}
	// release in steps
	time.Sleep(5 * time.Millisecond)
	c1.Close()
	time.Sleep(5 * time.Millisecond)
	c2.Close()
	for i := 0; i < waiters; i++ {
		<-done
	}
}

func TestIdleLifetimeEdgeBoundaries(t *testing.T) {
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) { return &fakeLDAP{}, nil }
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 1})
	ctx := context.Background()

	// set precise small durations
	p.maxIdleTime = 5 * time.Millisecond
	p.maxLifetime = 5 * time.Millisecond

	c, _ := p.Conn(ctx)
	c.Close()
	// exactly at edge; allow slight sleep less than edge to avoid flakiness
	time.Sleep(4 * time.Millisecond)
	c2, err := p.Conn(ctx)
	if err != nil {
		t.Fatalf("conn2: %v", err)
	}
	// may reuse same connection as not yet expired
	c2.Close()

	// now past edge
	time.Sleep(2 * time.Millisecond)
	c3, err := p.Conn(ctx)
	if err != nil {
		t.Fatalf("conn3: %v", err)
	}
	// must not be the same pooledConn as before
	if c3.conn == c2.conn {
		t.Fatalf("expected new connection after edge exceeded")
	}
	c3.Close()
}

func TestMetricsCountersIncrease(t *testing.T) {
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) { return &fakeLDAP{}, nil }
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 1})
	ctx := context.Background()

	// set very small times to trigger cleaner/fast-path closes
	p.maxIdleTime = 2 * time.Millisecond
	p.maxLifetime = 2 * time.Millisecond

	// cause waiting once
	c1, _ := p.Conn(ctx)
	done := make(chan struct{})
	go func() {
		c, _ := p.Conn(ctx)
		c.Close()
		done <- struct{}{}
	}()

	time.Sleep(1 * time.Millisecond)
	c1.Close()
	<-done

	// wait past expiration thresholds
	time.Sleep(10 * time.Millisecond)

	// force fast-path expiration via acquisition
	c2, err := p.Conn(ctx)
	if err != nil {
		t.Fatalf("conn after expiration: %v", err)
	}
	c2.Close()

	// read metrics
	waited := p.waitCount.Load()
	idleClosed := p.idleTimeClosedCount.Load()
	lifeClosed := p.lifeTimeClosedCount.Load()

	if waited == 0 {
		t.Fatalf("expected waitCount > 0")
	}
	if idleClosed == 0 && lifeClosed == 0 {
		t.Fatalf("expected some closed counters > 0 (idle or lifetime)")
	}
}

func TestDialBindErrorsDoNotExceedMaxOpen(t *testing.T) {
	// failing factory on first call, then success
	var calls int64
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) {
		if atomic.AddInt64(&calls, 1) == 1 {
			return nil, errors.New("dial fail")
		}
		return &fakeLDAP{}, nil
	}
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 1})
	ctx := context.Background()

	// first attempt fails, second should succeed and not exceed maxOpen
	_, _ = p.Conn(ctx)
	c, err := p.Conn(ctx)
	if err != nil {
		t.Fatalf("second conn should succeed: %v", err)
	}
	c.Close()
}

func TestConcurrentAcquireRelease(t *testing.T) {
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) {
		return &fakeLDAP{}, nil
	}
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 5})
	ctx := context.Background()

	const workers = 100
	const iters = 100
	done := make(chan struct{})

	for i := 0; i < workers; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for j := 0; j < iters; j++ {
				c, err := p.Conn(ctx)
				if err != nil {
					t.Errorf("conn: %v", err)
					return
				}
				time.Sleep(time.Millisecond)
				c.Close()
			}
		}()
	}
	for i := 0; i < workers; i++ {
		<-done
	}
	_ = p.Close()
}

func TestWaitersServedUnderContention(t *testing.T) {
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) {
		return &fakeLDAP{}, nil
	}
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 2})
	ctx := context.Background()

	c1, _ := p.Conn(ctx)
	c2, _ := p.Conn(ctx)

	start := make(chan struct{})
	const waiters = 20
	var got int32
	done := make(chan struct{})

	for i := 0; i < waiters; i++ {
		go func() {
			<-start
			c, err := p.Conn(ctx)
			if err == nil {
				atomic.AddInt32(&got, 1)
				c.Close()
			} else {
				t.Errorf("waiter conn err: %v", err)
			}
			done <- struct{}{}
		}()
	}

	close(start)
	time.Sleep(10 * time.Millisecond)
	c1.Close()
	time.Sleep(10 * time.Millisecond)
	c2.Close()

	for i := 0; i < waiters; i++ {
		<-done
	}
	if atomic.LoadInt32(&got) != waiters {
		t.Fatalf("not all waiters were served: got=%d want=%d", got, waiters)
	}
}

func TestCloseWhileWaiters(t *testing.T) {
	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) {
		return &fakeLDAP{}, nil
	}
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 1})
	ctx := context.Background()

	c, _ := p.Conn(ctx)

	const waiters = 10
	errs := make(chan error, waiters)
	for i := 0; i < waiters; i++ {
		go func() {
			_, err := p.Conn(ctx)
			errs <- err
		}()
	}

	_ = p.Close()
	c.Close()

	for i := 0; i < waiters; i++ {
		<-errs
	}
}

func TestNoLeakOnBadConnUnderConcurrency(t *testing.T) {
	type fakeLDAP2 struct{ fakeLDAP }
	var _ LdapConn = (*fakeLDAP2)(nil)

	factory := func(_ *LDAPAuthConfig, _ time.Duration) (LdapConn, error) {
		return &fakeLDAP2{}, nil
	}
	p := NewLDAPPool(LDAPPoolConfig{ConnFactory: factory, MaxConnections: 5})
	ctx := context.Background()

	const workers = 50
	const iters = 50
	done := make(chan struct{})

	for i := 0; i < workers; i++ {
		go func(i int) {
			defer func() { done <- struct{}{} }()
			for j := 0; j < iters; j++ {
				c, err := p.Conn(ctx)
				if err != nil {
					t.Errorf("conn: %v", err)
					return
				}
				if (i+j)%7 == 0 {
					le := &ldap.Error{ResultCode: ldap.ErrorNetwork, Err: errors.New("net")}
					c.conn.conn.(*fakeLDAP2).searchErr.Store(le)
					c.conn.conn.(*fakeLDAP2).hasErr.Store(true)
				} else {
					c.conn.conn.(*fakeLDAP2).hasErr.Store(false)
				}
				_, _ = c.Search(&ldap.SearchRequest{})
				c.Close()
			}
		}(i)
	}
	for i := 0; i < workers; i++ {
		<-done
	}
	_ = p.Close()
}
