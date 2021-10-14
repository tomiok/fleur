package main

import (
	"github.com/google/uuid"
	"github.com/ohchat-io/fleur/srv"
)

func main() {
	s := srv.NewServer("5566")
	cs := srv.ChatServer{
		ActiveConnections: make(map[string]*srv.Conn),
		TCPSrv:            s,
		Join:              make(chan *srv.Conn),
		Leave:             make(chan *srv.Conn),
		Input:             make(chan srv.Message),
	}
	go cs.Run()
	for {
		conn, _ := s.Listener.Accept()
		c := &srv.Conn{
			ID:         uuid.NewString(),
			Connection: conn,
		}
		go cs.HandleConnection(c)
	}
}
