package srv

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net"
	"sync"
	"time"
)

type ConnPool interface {
	Get() (net.Conn, error)
	GetWithTimeout(timeout time.Duration) (net.Conn, error)
	Close() error
	Remove(conn net.Conn) error
}

// pool implements ConnPool interface. Use channel buffer connections.
type pool struct {
	lock         sync.Mutex
	connections  chan *Conn
	minConnNum   int
	maxConnNum   int
	totalConnNum int
	closed       bool
	connCreator  func() (net.Conn, error)
}

var (
	errPoolIsClose = errors.New("connection pool has been closed")
	// Error for get connection time out.
	errTimeOut      = errors.New("get connection timeout")
	errContextClose = errors.New("get connection close by context")
)

// NewPool return new ConnPool. It bases on channel. It will start minConn connections in channel first.
// When Get()/GetWithTimeout called, if channel still has connection it will get connection from channel.
// Otherwise, pool check number of connection which had already created as the number are less than maxConn,
// it uses connCreator function to create new connection.
func NewPool(minConn, maxConn int) (*pool, error) {
	if minConn > maxConn || minConn < 0 || maxConn <= 0 {
		return nil, errors.New("number of connection bound error")
	}

	pool := &pool{}
	pool.minConnNum = minConn
	pool.maxConnNum = maxConn
	pool.connections = make(chan *Conn, maxConn)
	pool.closed = false
	pool.totalConnNum = 0
	err := pool.init()
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (p *pool) init() error {
	for i := 0; i < p.minConnNum; i++ {
		conn, err := p.createConn()
		if err != nil {
			return err
		}
		p.connections <- conn
	}
	return nil
}

// Get connection from connection pool. If connection poll is empty and already created connection number less than Max number of connection
// it will create new one. Otherwise, it wil wait someone put connection back.
func (p *pool) Get() (*Conn, error) {
	if p.isClosed() == true {
		return nil, errPoolIsClose
	}
	go func() {
		conn, err := p.createConn()
		if err != nil {
			return
		}
		p.connections <- conn
	}()
	select {
	case conn := <-p.connections:
		return p.packConn(conn), nil
	}
}

// GetWithTimeout can let you get connection wait for a time duration. If it cannot get connection in this time.
// It will return TimeOutError.
func (p *pool) GetWithTimeout(timeout time.Duration) (*Conn, error) {
	if p.isClosed() == true {
		return nil, errPoolIsClose
	}
	go func() {
		conn, err := p.createConn()
		if err != nil {
			return
		}
		p.connections <- conn
	}()
	select {
	case conn := <-p.connections:
		return p.packConn(conn), nil
	case <-time.After(timeout):
		return nil, errTimeOut
	}
}

func (p *pool) GetWithContext(ctx context.Context) (*Conn, error) {
	if p.isClosed() == true {
		return nil, errPoolIsClose
	}
	go func() {
		conn, err := p.createConn()
		if err != nil {
			return
		}
		p.connections <- conn
	}()
	select {
	case conn := <-p.connections:
		return p.packConn(conn), nil
	case <-ctx.Done():
		return nil, errContextClose
	}
}

// Close closes the connection pool. When close the connection pool it also close all connection already in connection pool.
// If connection not put back in connection it will not close. But it will close when it put back.
func (p *pool) Close() error {
	if p.isClosed() == true {
		return errPoolIsClose
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.closed = true
	close(p.connections)
	for conn := range p.connections {
		_ = conn.C.Close()
	}
	return nil
}

// Put can put connection back in connection pool. If connection has been closed, the connection will be close too.
func (p *pool) Put(conn *Conn) error {
	if p.isClosed() == true {
		return errPoolIsClose
	}
	if conn == nil {
		p.lock.Lock()
		p.totalConnNum = p.totalConnNum - 1
		p.lock.Unlock()
		return errors.New("cannot put nil to connection pool")
	}

	select {
	case p.connections <- conn:
		return nil
	default:
		return conn.C.Close()
	}
}

func (p *pool) isClosed() bool {
	p.lock.Lock()
	ret := p.closed
	p.lock.Unlock()
	return ret
}

// Remove let connection not belong connection pool.And it will close connection.
func (p *pool) Remove(c *Conn) error {
	if p.isClosed() == true {
		return errPoolIsClose
	}

	p.lock.Lock()
	p.totalConnNum = p.totalConnNum - 1
	p.lock.Unlock()
	return c.C.Close()
}

// createConn will create one connection from connCreator. And increase connection counter.
func (p *pool) createConn() (*Conn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.totalConnNum >= p.maxConnNum {
		return nil, fmt.Errorf("connot Create new connection. Now has %d.Max is %d", p.totalConnNum, p.maxConnNum)
	}
	conn, err := p.connCreator()
	if err != nil {
		return nil, fmt.Errorf("cannot create new connection.%s", err)
	}
	p.totalConnNum = p.totalConnNum + 1
	return &Conn{
		ID: uuid.NewString(),
		C:  conn,
	}, nil
}

func (p *pool) packConn(conn *Conn) *Conn {
	ret := &Conn{p: p}
	ret.C = conn.C
	return ret
}

type Conn struct {
	ID string
	C  net.Conn
	p  *pool
}
