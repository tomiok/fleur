package srv

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Hub struct {
}

type TCPServer struct {
	addr string
	L    net.Listener
}

func NewServer(addr string) *TCPServer {
	l, err := net.Listen("tcp4", ":"+addr)
	if err != nil {
		panic(err.Error())
	}
	return &TCPServer{
		addr: addr,
		L:    l,
	}
}

func HandleConnection(c *Conn) {
	for {
		msg, err := bufio.NewReader(c.C).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(msg)

		if temp == "STOP" {
			break
		}
		fmt.Println(c.ID + "> " + temp)

	}
	_ = c.C.Close()
}
