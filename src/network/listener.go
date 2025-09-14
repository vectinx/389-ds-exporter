package network

import (
	"net"
	"time"
)

type TimeoutListener struct {
	net.Listener
	readTimeout time.Duration
}

func (tl *TimeoutListener) Accept() (net.Conn, error) {
	conn, err := tl.Listener.Accept()
	if err != nil {
		return nil, err
	}

	// Устанавливаем таймаут чтения (например, чтобы избежать зависших telnet-соединений)
	if err := conn.SetReadDeadline(time.Now().Add(tl.readTimeout)); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}

func NewTimeoutListener(inner net.Listener, readTimeout time.Duration) *TimeoutListener {
	return &TimeoutListener{
		Listener:    inner,
		readTimeout: readTimeout,
	}
}
