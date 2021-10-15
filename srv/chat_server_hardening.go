package srv

import "github.com/rs/zerolog/log"

func (server *ChatServer) CloseConnection(c *Conn) {
	delete(server.ActiveConnections, c.Nick)
	err := c.Connection.Close()

	if err != nil {
		log.Warn().Msgf("cannot clone connection %s", err.Error())
	}

	server.Leave <- c
}
