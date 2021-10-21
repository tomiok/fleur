package srv

import "net"

type Conn struct {
	Nick       string
	Connection net.Conn
	Wait       chan struct{}
}

func NewConn(connection net.Conn) *Conn {
	return &Conn{
		Connection: connection,
		Wait:       make(chan struct{}),
	}
}
