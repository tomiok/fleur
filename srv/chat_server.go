package srv

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"strings"
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

func NewChatServer(port string) *ChatServer {
	return &ChatServer{
		ActiveConnections: make(map[string]*Conn),
		TCPSrv:            NewServer(port),
		Join:              make(chan *Conn),
		Leave:             make(chan *Conn),
		Input:             make(chan Message),
	}
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

func (server *ChatServer) HandleConnection(c *Conn) {
	defer func() {
		//TODO handle this by deleting one connection in the server as well
		_ = c.Connection.Close()
	}()
	for {
		Write(c.Connection, "Enter your nick:")

		scanner := bufio.NewScanner(c.Connection)
		scanner.Scan()
		c.Nick = scanner.Text()
		server.Join <- c

		// Read and write the message. Lookup the receiver.
		go func() {
			for scanner.Scan() {
				ln := scanner.Text()
				s := strings.Split(ln, " ")
				//TODO handle when a user repeat the nickname
				//TODO handle bad messages
				//TODO format messages correctly
				user := server.ActiveConnections[s[0]]
				Write(user.Connection, ln)
			}
		}()

		// wait for it
		<-c.Wait
	}
}

func Write(w io.Writer, msg string) {
	_, err := io.WriteString(w, msg)
	if err != nil {
		log.Warn().Msgf("cannot write message %s", err.Error())
	}
}
