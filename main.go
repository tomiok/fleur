package main

import (
	"fmt"
	"github.com/ohchat-io/fleur/srv"
	"net"
)

const minCount = 3

func main() {

	pool, err := srv.NewPool(minCount, minCount*2, func() (net.Conn, error) {
		return net.Dial("tcp", "127.0.0.1:")
	})

	if err != nil {
		panic(err.Error())
	}

	s := srv.NewServer("5566")

	for {
		conn, err := pool.Get()
		if err != nil {
			fmt.Println("cannot get conn " + err.Error())
		}
		go s.HandleConnection(conn.C)
	}
}
