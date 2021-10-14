package srv

import "fmt"

func (server *ChatServer) CloseConnection(c *Conn) {
	fmt.Println("closing")
	server.Leave <- c
}
