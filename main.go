package main

import (
	"log"
	"net"
	"sync"
	"time"
)

const (
	defaultAddr = "stun.l.google.com:19302"
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
		Connection: conn,
	})
}

func NewClient(opt ClientOptions) (*Client, error) {
	client := &Client{}
	return client, nil
}
