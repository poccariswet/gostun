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
	close       chan struct{}
	agent       messageClient
}

type messageClient interface {
	Process(*Message) error
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

	c.wg.Add(1)
	go c.readUntil()

	return c, nil
}

func (c *Client) readUntil() {
	defer c.wg.Done()

	m := new(Message)
	m.Raw = make([]byte, 1024)
	for {
		select {
		case <-c.close:
			return
		default:
		}
		_, err := m.ReadConn(c.conn)
		if err == nil {
			if Err := c.agent.Process(m); Err != nil {
				return
			}
		}
	}
}
