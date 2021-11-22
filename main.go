package main

import (
	"github.com/google/uuid"
	"github.com/ohchat-io/fleur/srv"
)

func main() {
	s := srv.NewServer("5566")

	for {
		conn, _ := s.L.Accept()
		c := &srv.Conn{
			ID: uuid.NewString(),
			C:  conn,
		}
		go srv.HandleConnection(c)
	}
}
