package srv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
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
	// message properties
	Type       string `json:"type"`     // define which type of message is, and who is going to receive
	Sender     string `json:"sender"`   // the sender
	Receiver   string `json:"receiver"` // the receiver
	Body       string `json:"body"`     // the body
	ExcludeOne string `json:"-"`        // if a connection should be excluded for the message
	// common properties
	Connections []string `json:"connections"` //the people connected.
}

func (m *Message) Build() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Error().Msgf("cannot marshall msg %s", err.Error())
		return "", err
	}

	return string(b), nil
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
					Type:       msgTypeBroadcast,
					Sender:     systemSender,
					ExcludeOne: conn.Nick,
					Body:       fmt.Sprintf("%s joined Fleur channel", conn.Nick),
				}
			}()
		// When a user leaves the server, send a message to everyone.
		case conn := <-server.Leave:
			go func() {
				server.Input <- Message{
					Type:       msgTypeBroadcast,
					Sender:     systemSender,
					ExcludeOne: conn.Nick,
					Body:       fmt.Sprintf("%s left Fleur channel", conn.Nick),
				}
			}()
		case msg := <-server.Input:
			var t = msg.Type
			switch t {
			case msgTypeBroadcast:
				for k, v := range server.ActiveConnections {
					if k != msg.ExcludeOne {
						WriteMessage(v.Connection, msg.Build)
					}
				}
			case msgTypeDirect:
				receiver := server.GetConnection(msg.Receiver)
				// prevent send message to a non-connected user
				if receiver != nil {
					WriteMessage(receiver.Connection, msg.Build)
				}
			case msgTypeSelf:
				receiver := server.GetConnection(msg.Receiver)
				if receiver != nil {
					WriteMessage(receiver.Connection, msg.Build)
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
			nick := scanner.Text()
			c.Nick = nick
			if !server.IsValidNickname(nick) {
				break
			}
		}

		// Emit a new join event.
		server.Join <- c

		server.Input <- Message{
			Type:        msgTypeSelf,
			Sender:      systemSender,
			Receiver:    c.Nick,
			Body:        "welcome " + c.Nick,
			Connections: server.GetConnections(),
		}

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

func WriteMessage(w io.Writer, f func() (string, error)) {
	msg, err := f()
	if err != nil {
		return
	}

	write(w, msg)
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
