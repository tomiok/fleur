package main

import (
	"github.com/ohchat-io/fleur/config"
	"github.com/ohchat-io/fleur/srv"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	conf := config.Bind()
	conf.PrintConfigs()

	cs := srv.NewChatServer(conf.Port)

	go cs.Run()

	for {
		conn, err := cs.TCPSrv.Listener.Accept()
		if err != nil {
			log.Error().Msgf("cannot accept connections. %s", err.Error())
			os.Exit(1)
		}
		go cs.HandleConnection(srv.NewConn(conn))
	}
}
