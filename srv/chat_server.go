package srv

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"strings"
)

const (
	broadcast = "BROADCAST"
)

type ChatServer struct {
	ActiveConnections map[string]*Conn
	TCPSrv            *TCPServer
	Join              chan *Conn
	Leave             chan *Conn
	Input             chan Message
}

type Message struct {
	Type     string
	Sender   string
	Receiver string
	Text     string
}

type Conn struct {
	Nick       string
	Connection net.Conn
	Wait       chan struct{}
}

func NewConn(connection net.Conn) *Conn {
	return &Conn{Connection: connection}
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

// Run starts the server
func (server *ChatServer) Run() {
	for {
		select {
		case conn := <-server.Join:
			server.AddUser(conn)
			go func() {
				server.Input <- Message{
					Type:   broadcast,
					Sender: conn.Nick,
					Text:   fmt.Sprintf("%s joined Fleur channel", conn.Nick),
				}
			}()
		case msg := <-server.Input:
			var t = msg.Type
			switch t {
			case broadcast:
				for k, v := range server.ActiveConnections {
					if k != msg.Sender {
						Write(v.Connection, msg.Text)
					}
				}
			}
		case conn := <-server.Leave:
			go func() {
				delete(server.ActiveConnections, conn.Nick)
				server.Input <- Message{
					Type:   broadcast,
					Sender: conn.Nick,
					Text:   fmt.Sprintf("%s left Fleur channel", conn.Nick),
				}
			}()
			func() {
				err := conn.Connection.Close()
				if err != nil {
					log.Warn().Msgf("cannot clone connection %s", err.Error())
				}
			}()
		}
	}
}

func (server *ChatServer) HandleConnection(c *Conn) {

	for {

		Write(c.Connection, "Enter your nick: ")

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
	_, err := io.WriteString(w, fmt.Sprintf("%s \n", msg))
	if err != nil {
		log.Warn().Msgf("cannot write message %s", err.Error())
	}
}
