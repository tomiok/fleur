package main

import (
	"github.com/ohchat-io/fleur/config"
	"github.com/ohchat-io/fleur/srv"
)

func main() {
	conf := config.Bind()
	conf.PrintConfigs()

	cs := srv.NewChatServer(conf.Port)

	go cs.Run()

	for {
		conn, _ := cs.TCPSrv.Listener.Accept()
		go cs.HandleConnection(srv.NewConn(conn))
	}
}
