package network

import (
	"net"
	"time"
)

// TimeoutListener is the wrapper over net.Listener,
// implementing timeout for sending the first byte of tcp connections.
//
// Needed to prevent the appearance of frozen connections - open tcp - sessions that do not send data
type TimeoutListener struct {
	net.Listener
	readTimeout time.Duration
}

// Accept waits for and returns the next connection to the listener with first byte read timeout.
func (tl *TimeoutListener) Accept() (net.Conn, error) {
	conn, err := tl.Listener.Accept()
	if err != nil {
		return nil, err
	}

	// Set a read timeout (for example, to avoid hanging telnet connections)
	if err := conn.SetReadDeadline(time.Now().Add(tl.readTimeout)); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}

// NewTimeoutListener function creates new TimeoutListener with specified timeout
func NewTimeoutListener(inner net.Listener, readTimeout time.Duration) *TimeoutListener {
	return &TimeoutListener{
		Listener:    inner,
		readTimeout: readTimeout,
	}
}
