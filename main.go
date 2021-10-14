package main

import (
	"github.com/google/uuid"
	"github.com/ohchat-io/fleur/config"
	"github.com/ohchat-io/fleur/srv"
)

//TODO fix broadcasting, do a correct one
//TODO add init functions for TCP and chat servers
func main() {
	conf := config.Bind()
	conf.PrintConfigs()
	s := srv.NewServer(conf.Port)
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
