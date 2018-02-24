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

	if c.agent == nil {
		c.agent = NewAgent()
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
		_, err := m.ReadConn(c.conn) // read and decode message
		if err == nil {
			if Err := c.agent.Process(m); Err != nil {
				return
			}
		}
	}
}

// process of transaction in message
type Agent struct {
	transactions map[transactionID]TransactionAgent
	mux          sync.Mutex
	nonHandler   Handler // non-registered transactions
}

type transactionID [TransactionIDSize]byte

// transaction in progress
type TransactionAgent struct {
	ID      transactionID
	Timeout time.Time
}

type AgentHandle struct {
	handler Handler
}

type Handler interface {
	HandleEvent(e EventObject)
}

func NewAgent() *Agent {
	h := AgentHandle{}
	a := &Agent{
		transactions: make(map[transactionID]TransactionAgent),
		nonHandler:   h.handler,
	}
	return a
}
