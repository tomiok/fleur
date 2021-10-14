package srv

import (
	"fmt"
	"net"
)

type ChatServer struct {
	ActiveConnections map[string]*Conn
	TCPSrv            *TCPServer
	Join              chan *Conn
	Leave             chan *Conn
	Input             chan Message
}

type Message struct {
	Receiver string
	Text     string
}

type Conn struct {
	ID         string
	Nick       string
	Connection net.Conn
	Wait       chan struct{}
}

func (server *ChatServer) AddUser(c *Conn) {
	server.ActiveConnections[c.Nick] = c
}

func (server *ChatServer) Run() {
	for {
		select {
		case u := <-server.Join:
			server.AddUser(u)
			go func() {
				server.Input <- Message{
					Receiver: "BROADCAST",
					Text:     u.Nick + " joined",
				}
			}()
		case msg := <-server.Input:
			var r = msg.Receiver
			switch r {
			case "BROADCAST":
				fmt.Println("broadcast")
			default:

			}


		}
	}
}
