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
		io.WriteString(c.Connection, "Enter your nick:")

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
				user := server.ActiveConnections[s[0]]
				io.WriteString(user.Connection, ln)
			}
		}()

	}
}

func Write(w io.Writer, msg string) {
	_, err := io.WriteString(w, msg)
	if err != nil {
		log.Warn().Msgf("cannot write message %s", err.Error())
	}
}
