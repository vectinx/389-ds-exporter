package ldap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUrlParseError(t *testing.T) {
	auth := &LDAPAuthConfig{URL: ":::bad-url"}

	_, err := RealConnectionDialUrl(auth)
	require.Error(t, err, "Passing invalid URL should fail")
	if err == nil {
		t.Fatal("")
	}
}

func TestDialFail(t *testing.T) {
	auth := &LDAPAuthConfig{
		URL:         "ldap://127.0.0.1:389",
		DialTimeout: 1 * time.Second,
	}

	_, err := RealConnectionDialUrl(auth)
	require.Error(t, err, "Connecting to a non-existent server should fail with an error")
}
