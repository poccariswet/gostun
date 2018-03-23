package gostun

import (
	"io"
	"log"
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
	rw          sync.RWMutex
	clientclose bool
}

type messageClient interface {
	ProcessHandle(*Message) error
	TimeOutHandle(time.Time) error
	TransactionHandle([TransactionIDSize]byte, Handler, time.Time) error
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
		agent:       NewAgent(),
		TimeoutRate: defaultTimeoutRate,
	}

	c.wg.Add(2)
	go c.readDecode()
	go c.timeoutUntil()

	return c, nil
}

func (c *Client) readDecode() {
	defer c.wg.Done()

	m := new(Message)
	m.Raw = make([]byte, 1024)

	_, err := m.ReadConn(c.conn) // read and decode message
	if err == nil {
		if processErr := c.agent.ProcessHandle(m); processErr == ErrAgent {
			return
		}
	} else {
		log.Print(err)
	}
}

func (c *Client) timeoutUntil() {
	t := time.NewTicker(c.TimeoutRate)
	defer c.wg.Done()
	for {
		select {
		case <-c.close:
			t.Stop()
			return
		case trate := <-t.C:
			err := c.agent.TimeOutHandle(trate)
			if err == nil || err == ErrAgent {
				return
			}
			panic(err)
		}
	}
}
