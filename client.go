package gostun

import (
	"errors"
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
	Collect(time.Time) error
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

	c.wg.Add(2)
	go c.readUntil()
	go c.collectUntil()

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
			if processErr := c.agent.Process(m); processErr != nil {
				return
			}
		}
	}
}

var ErrAgent = errors.New("agent closed")

func (c *Client) collectUntil() {
	t := time.NewTicker(c.TimeoutRate)
	defer c.wg.Done()
	for {
		select {
		case <-c.close:
			t.Stop()
			return
		case trate := <-t.C:
			err := c.agent.Collect(trate)
			if err == nil || err == ErrAgent {
				return
			}
			panic(err)
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
	handler Handler // if transaction is succeed will be called
}

type AgentHandle struct {
	handler Handler
}

// reference http.HandlerFunc same work
type Handler interface {
	HandleEvent(e EventObject)
}

type HandleFunc func(e EventObject)

func (f HandleFunc) HandleEvent(e EventObject) {
	f(e)
}

type EventObject struct {
	Msg *Message
	err error
}

func NewAgent() *Agent {
	h := AgentHandle{}
	a := &Agent{
		transactions: make(map[transactionID]TransactionAgent),
		nonHandler:   h.handler,
	}
	return a
}

func (a *Agent) Process(m *Message) error {
	e := EventObject{
		Msg: m,
	}
	a.mux.Lock() // protect transaction
	tr, ok := a.transactions[m.TransactionID]
	delete(a.transactions, m.TransactionID) //delete maps entry

	if ok {
		tr.handler.HandleEvent(e) // HandleEvent cast the e to hander type
	} else if a.nonHandler != nil {
		a.nonHandler.HandleEvent(e) // the transaction is not registered
	}
	return nil
}
