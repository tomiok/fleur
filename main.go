package main

import (
	"github.com/ohchat-io/fleur/srv"
)

const minCount = 3

func main() {

	s := srv.NewServer("5566")

	for {
		conn, _ := s.L.Accept()
		go srv.HandleConnection(conn)
	}
}
