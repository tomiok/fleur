package srv

import (
	"io"
	"net"
	"time"
)

const (
	MaxReadLength    = 1024
	ReadWriteTimeout = time.Minute
)

type Conn struct {
	Nick          string
	Connection    net.Conn
	Wait          chan struct{}
	Limiter       *io.LimitedReader
	MaxReadLength int64
	ReadWriteTimeout time.Duration
}

func NewConn(conn net.Conn) *Conn {
	_ = conn.SetReadDeadline(time.Now().Add(ReadWriteTimeout))
	limiter := &io.LimitedReader{
		R: conn,
		N: MaxReadLength,
	}
	return &Conn{
		Connection:    conn,
		Wait:          make(chan struct{}),
		Limiter:       limiter,
		MaxReadLength: MaxReadLength,
		ReadWriteTimeout: ReadWriteTimeout,
	}
}
