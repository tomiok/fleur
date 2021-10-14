package srv

import (
	"log"
	"net"
)

type TCPServer struct {
	addr     string
	Listener net.Listener
}

func NewServer(addr string) *TCPServer {
	l, err := net.Listen("tcp4", ":"+addr)

	if err != nil {
		log.Fatalf("cannot create listener for tcp4 connection, %s", err.Error())
		return nil
	}

	return &TCPServer{
		addr:     addr,
		Listener: l,
	}
}
