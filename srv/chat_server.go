package srv

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net"
)

// ChatServer is the main structure that holds all the necessary information for the tcp and web server
// TODO add stats for connections and messages
type ChatServer struct {
	ActiveConnections map[string]*Conn
	TCPSrv            *TCPServer
	Join              chan *Conn
	Leave             chan *Conn
	Input             chan Message
}

type Message struct {
	Type        string
	Sender      string
	Receiver    string
	Text        string
	Connections []string
}

type Conn struct {
	Nick       string
	Connection net.Conn
	Wait       chan struct{}
}

func NewConn(connection net.Conn) *Conn {
	return &Conn{
		Connection: connection,
		Wait:       make(chan struct{}),
	}
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

// Run starts the server, keep all the channels running and listening for new events.
// Join, Leave and Input, those are the main events in the chat server.
func (server *ChatServer) Run() {
	for {
		select {
		// When a user joins the server, send a message to everyone.
		case conn := <-server.Join:
			server.AddUser(conn)
			go func() {
				server.Input <- Message{
					Type:   msgTypeBroadcast,
					Sender: conn.Nick,
					Text:   fmt.Sprintf("%s joined Fleur channel", conn.Nick),
				}
			}()
		// When a user leaves the server, send a message to everyone.
		case conn := <-server.Leave:
			go func() {
				server.Input <- Message{
					Type:   msgTypeBroadcast,
					Sender: conn.Nick,
					Text:   fmt.Sprintf("%s left Fleur channel", conn.Nick),
				}
			}()
		case msg := <-server.Input:
			var t = msg.Type
			switch t {
			case msgTypeBroadcast:
				for k, v := range server.ActiveConnections {
					if k != msg.Sender {
						WriteMessage(v.Connection, "", msg.Text)
					}
				}
			case msgTypeDirect:
				receiver := server.ActiveConnections[msg.Receiver]
				// prevent send message to a non-connected user
				if receiver != nil {
					WriteMessage(receiver.Connection, msg.Sender, msg.Text)
				}

			}
		}
	}
}

func (server *ChatServer) HandleConnection(c *Conn) {
	for {
		scanner := bufio.NewScanner(c.Connection)
		for {
			WritePrompt(c.Connection, "Enter your nick: ")
			scanner.Scan()
			c.Nick = scanner.Text()
			if !server.IsValidNickname(c.Nick) {
				break
			}
			// TODO send all the information about the server to the current connection
		}

		// Emit a new join event.
		server.Join <- c

		// Read and write the message. Lookup the receiver.
		go func() {
			for scanner.Scan() {
				text := scanner.Text()
				msg, err := parse(c.Nick, text, directMsgParser)
				if err != nil {
					continue
				}

				// write to receiver(s)
				server.Input <- msg
			}
			// tidy up
			server.CloseConnection(c)
		}()

		// Wait for it.
		<-c.Wait
	}
}

func WriteMessage(w io.Writer, sender, msg string) {
	write(w, sender+"> "+msg+"\n")
}

func WritePrompt(w io.Writer, msg string) {
	write(w, msg)
}

func write(w io.Writer, msg string) {
	_, err := io.WriteString(w, fmt.Sprintf("%s", msg))
	if err != nil {
		log.Warn().Msgf("cannot write message %s", err.Error())
	}
}
