package srv

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

type Hub struct {
	*pool
}

type TCPServer struct {
	addr string
	L    net.Listener
}

func NewServer(addr string) *TCPServer {
	l, err := net.Listen("tcp4", ":"+addr)
	if err != nil {
		panic(err.Error())
	}
	return &TCPServer{
		addr: addr,
		L:    l,
	}
}

func HandleConnection(c *Conn) {
	for {
		netData, err := bufio.NewReader(c.C).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(netData)
		if temp == "STOP" {
			break
		}
		fmt.Println(c.ID + "> " + temp)

	}
	_ = c.C.Close()
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

func Write(s net.Conn, data []byte) error {
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	_, err := s.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func Read(s net.Conn) ([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(s, header)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(s, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
