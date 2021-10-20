package srv

import (
	"errors"
	"github.com/rs/zerolog/log"
	"strings"
)

func (server *ChatServer) CloseConnection(c *Conn) {
	delete(server.ActiveConnections, c.Nick)
	err := c.Connection.Close()

	if err != nil {
		log.Warn().Msgf("cannot clone connection %s", err.Error())
	}

	server.Leave <- c
}

func parse(sender, msg string, parser func(s, m string) (Message, error)) (Message, error) {
	return parser(sender, msg)
}

func directMsgParser(sender, message string) (Message, error) {
	if message == "" {
		return Message{}, errors.New("empty message")
	}

	values := strings.Split(message, " ")

	receiver := values[0]
	text := strings.Join(values[1:], " ")

	return Message{
		Sender:   sender,
		Receiver: receiver,
		Text:     text,
		Type:     msgTypeDirect,
	}, nil
}

func (server *ChatServer) IsValidNickname(nick string) bool {
	_, b := server.ActiveConnections[nick]
	return b
}

func (server *ChatServer) ShowConnections() []string {
	actives := server.ActiveConnections
	connections := make([]string, len(actives))

	for k := range actives {
		connections = append(connections, k)
	}

	return connections
}

func (server *ChatServer) GetConnection(nick string) *Conn{
	return server.ActiveConnections[nick]
}