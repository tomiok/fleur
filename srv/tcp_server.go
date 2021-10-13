package srv

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Hub struct {
	*pool
}

type TCPServer struct {
	addr     string
	listener net.Listener
}

func NewServer(addr string) *TCPServer {
	l, err := net.Listen("tcp4", ":"+addr)
	if err != nil {
		panic(err.Error())
	}
	return &TCPServer{
		addr:     addr,
		listener: l,
	}
}

func (s *TCPServer) HandleConnection(c net.Conn) {
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(netData)
		if temp == "STOP" {
			break
		}
		fmt.Println(temp)

	}
	_ = c.Close()
}

func (h *Hub) Destroy(conn *Conn) error {
	if h.pool == nil {
		return errors.New("connection not belong any connection pool")
	}
	err := h.pool.Remove(conn)
	if err != nil {
		return err
	}
	h.pool = nil
	return nil
}

// Close will push connection back to connection pool. It will not close the real connection.
func (h *Hub) Close(conn *Conn) error {
	if h.pool == nil {
		return errors.New("connection not belong any connection pool")
	}
	return h.pool.Put(conn)
}
