package srv

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"strings"
)

type TCPServer struct {
	addr     string
	Listener net.Listener
}

func NewServer(addr string) *TCPServer {
	l, err := net.Listen("tcp4", ":"+addr)
	if err != nil {
		panic(err.Error())
	}
	return &TCPServer{
		addr:     addr,
		Listener: l,
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
