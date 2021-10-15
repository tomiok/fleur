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

func (server *ChatServer) MessageParser(sender, message string) (Message, error) {
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
	}, nil
}
