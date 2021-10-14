package main

import (
	"github.com/google/uuid"
	"github.com/ohchat-io/fleur/config"
	"github.com/ohchat-io/fleur/srv"
)

//TODO fix broadcasting, do a correct one
func main() {
	conf := config.Bind()
	conf.PrintConfigs()

	cs := srv.NewChatServer(conf.Port)

	go cs.Run()

	for {
		conn, _ := cs.TCPSrv.Listener.Accept()
		c := &srv.Conn{
			ID:         uuid.NewString(),
			Connection: conn,
		}
		go cs.HandleConnection(c)
	}
}
