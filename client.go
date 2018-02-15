//package gostun
package gostun

import (
	"io"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn        Connection
	TimeoutRate time.Duration
	wg          sync.WaitGroup
}

type Connection interface {
	io.Reader
	io.Writer
	io.Closer
}

const defaultTimeoutRate = time.Millisecond * 100

func Dial(network, addr string) (*Client, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn)
}

func NewClient(conn net.Conn) (*Client, error) {
	c := &Client{
		conn:        conn,
		TimeoutRate: defaultTimeoutRate,
	}

	return c, nil
}
