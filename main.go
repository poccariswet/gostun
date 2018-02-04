package main

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	defaultAddr         = "stun.l.google.com:19302"
	TransactionIDSize   = 12 // 96 bit
	defaultTime         = time.Millisecond * 100
	stunMessageHeader   = 20
	magicCookie         = 0x2112A442
	stunAttributeHeader = 4
)

type Client struct {
	agent     ClientAgent
	conn      Connection
	close     chan struct{}
	closed    bool
	closedMux sync.RWMutex
	gcRate    time.Duration
	wg        sync.WaitGroup
}

type ClientOptions struct {
	Agent      ClientAgent
	Connection Connection
	TimeOut    time.Duration
}

type ClientAgent interface {
	Process(*Message) error
	Close() error
	Start(id [TransactionIDSize]byte, deadline time.Time, f Handler) error
	Stop(id [TransactionIDSize]byte) error
	Collect(time.Time) error
}

type Connection interface {
	io.Reader
	io.Writer
	io.Closer
}

type Agent struct {
	transactions map[transactionID]agentTransaction
	closed       bool
	mux          sync.Mutex
	zeroHandler  Handler
}

type transactionID [TransactionIDSize]byte

type AgentOptions struct {
	Handler Handler
}

type Handler interface {
	HandleEvent(e Event)
}

type Message struct {
	Type          MessageType
	Length        uint32
	TransactionID [TransactionIDSize]byte
	Attributes    Attributes
	Raw           []byte
}

type MessageType struct {
	Method Method       // binding
	Class  MessageClass // request
}

type MessageClass byte

func main() {
	c, err := Dial("udp", defaultAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func Dial(network, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return NewClient(ClientOptions{
		Agent:      NewAgent(AgentOptions{}),
		Connection: conn,
		TimeOut:    defaultTime,
	})
}

func NewClient(opt ClientOptions) (*Client, error) {
	client := &Client{
		agent:  opt.Agent,
		conn:   opt.Connection,
		gcRate: opt.TimeoutRate,
	}

	client.wg.Add(2)
	go client.readUntilClosed()
	go client.collectUntilClosed()
	return client, nil
}

func NewAgent(Aopt AgentOptions) *Agent {
	return &Agent{
		transactions: make(map[transactionID]agentTransaction),
		zeroHandler:  Aopt.Handler,
	}
}

func (c *Client) readUntilClosed() {
	defer c.wg.Done()

	m := new(Message)
	m.Raw = make([]byte, 1024) // Message struct のRawに1024バイトの空きを作る

	for {
		select {
		case <-c.close():
			return
		default:
		}
		_, err := m.ReadFrom(c.conn) // m.Raw
		if err == nil {
			if Err := c.agent.Process(m); Err == ErrAgentClosed {
				return
			}
		}
	}
}

func (m *Message) ReadFrom(r io.Reader) (int64, error) {
	tBuf := m.Raw[:cap(m.Raw)]

	n, err := r.Read(tBuf)
	if err != nil {
		return int64(n), err
	}
	m.Raw = tBuf[:n]
	return int64(n), m.Decode()
}
