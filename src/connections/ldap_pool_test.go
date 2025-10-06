package connections

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/require"
)

type mockLdapConn struct {
	BindFn   func(string, string) error
	SearchFn func(*ldap.SearchRequest) (*ldap.SearchResult, error)
	UnbindFn func() error

	BindCount   atomic.Int32
	SearchCount atomic.Int32
	UnbindCount atomic.Int32
}

func (c *mockLdapConn) Bind(bind_dn string, bind_pw string) error {
	c.BindCount.Add(1)
	if c.BindFn == nil {
		return nil
	}
	return c.BindFn(bind_dn, bind_pw)
}

func (c *mockLdapConn) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	c.SearchCount.Add(1)
	if c.SearchFn == nil {
		return &ldap.SearchResult{}, nil
	}
	return c.SearchFn(req)
}
func (c *mockLdapConn) Unbind() error {
	c.UnbindCount.Add(1)
	if c.UnbindFn == nil {
		return nil
	}
	return c.UnbindFn()
}

// Test CSN Structure.
func TestLdapPool(t *testing.T) {
	t.Log("Given the needed to check the operation of the LdapConnectionPool")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen we get and return a connection from the newly created pool", testID)
		{
			mockConn := mockLdapConn{
				SearchFn: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
					return &ldap.SearchResult{}, nil
				},
			}
			ldapConnPoolConfig := LdapConnectionPoolConfig{
				ServerURL:      "ldap://localhost:389",
				BindDN:         "cn=directory manager",
				BindPw:         "12345678",
				MaxConnections: 1,
				ConnFactory: func(string, context.Context, time.Duration) (LdapConn, error) {
					return &mockConn, nil
				},
			}

			ldapConnPool := NewLdapConnectionPool(ldapConnPoolConfig)

			t.Logf("\t\tCreate and return connection")
			{
				get_ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				defer cancel()
				conn, err := ldapConnPool.Get(get_ctx)
				require.NoError(t, err, "Connections from a non-empty pool must be created without errors")

				// Trying to get more connections than there are in the pool.
				_, err = ldapConnPool.Get(get_ctx)
				require.ErrorIs(
					t,
					err,
					ErrPoolGetTimedOut,
					"An attempt to obtain a connection from a pool that has exceeded the limit should return an error after a timeout",
				)

				// Closing connection
				conn.Close()
				// Closing the connection again should cause panic
				require.Panics(t, func() { conn.Close() })

				require.Equal(t, int32(1), mockConn.BindCount.Load(), "A connection must have one BIND operation after use")
			}

			t.Logf("\t\tGet a previously created connection from the pool")
			{
				get_ctx, cancel := context.WithTimeout(
					context.Background(),
					500*time.Millisecond,
				)

				defer cancel()
				conn, err := ldapConnPool.Get(get_ctx)
				require.NoError(t, err, "Getting a connection previously created in the pool should not fail")

				require.Equal(
					t,
					int32(1),
					mockConn.SearchCount.Load(),
					"The reissued connection must have at least one search performed - connection check",
				)
				conn.Close()

				require.Equal(
					t,
					1,
					ldapConnPool.ConnsAtPool(),
					"After returning a connection, the pool should contain as many connections as were previously allocated",
				)
			}

			t.Logf("\t\tClose pool")
			{
				close_ctx, close_cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer close_cancel()

				_ = ldapConnPool.Close(close_ctx)
				require.Equal(t, 0, ldapConnPool.ConnsAtPool(), "After closing the pool, there should be no connections left in it")

				require.Equal(
					t,
					int32(1),
					mockConn.UnbindCount.Load(),
					"A connection must have one UNBIND operation after pool closing",
				)
			}

			t.Logf("\t\tTest the operation of the closed pool")
			{
				get_ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				defer cancel()
				close_ctx, close_cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer close_cancel()
				err := ldapConnPool.Close(close_ctx)
				require.ErrorIs(t, err, ErrPoolClosed, "Closing the pool again should fail with 'ErrPoolClosed'")

				_, err = ldapConnPool.Get(get_ctx)
				require.ErrorIs(
					t,
					err,
					ErrPoolClosed,
					"Attempting to obtain a connection from the closed pool should fail with 'ErrPoolClosed'",
				)
			}

		}

		testID++
		t.Logf("\tTest %d:\tWhen competitive access to a pool occurs", testID)
		{
			mockConn := mockLdapConn{
				BindFn: func(string, string) error {
					// We simulate a small delay in the LDAP request to ensure that pool connection requests are queued.
					time.Sleep(50 * time.Millisecond)
					return nil
				},
				SearchFn: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
					// We simulate a small delay in the LDAP request to ensure that pool connection requests are queued.
					time.Sleep(50 * time.Millisecond)
					return &ldap.SearchResult{}, nil
				},
			}
			ldapConnPoolConfig := LdapConnectionPoolConfig{
				ServerURL:      "ldap://localhost:389",
				BindDN:         "cn=directory manager",
				BindPw:         "12345678",
				MaxConnections: 10,
				ConnFactory: func(string, context.Context, time.Duration) (LdapConn, error) {
					return &mockConn, nil
				},
			}

			ldapConnPool := NewLdapConnectionPool(ldapConnPoolConfig)

			// This test itself doesn't check anything and is needed to detect concurrent execution errors.
			// For example, using the `go test -race`
			t.Logf("\t\tCompetitively obtaining from a pool of connections in quantities greater than the pool size")
			{
				var wg sync.WaitGroup
				for range 1000 {
					wg.Add(1)
					go func() {
						defer wg.Done()

						get_ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
						defer cancel()

						conn, _ := ldapConnPool.Get(get_ctx)
						if conn != nil {
							defer conn.Close()
							_, _ = conn.Search(&ldap.SearchRequest{})
						}
					}()
				}

				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()
				_ = ldapConnPool.Close(ctx)
				wg.Wait()

			}
		}
	}
}
