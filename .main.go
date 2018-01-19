package main

import (
	"fmt"
	"log"
	"net"
)

const (
	DefaultName   = "stun_client"
	DefaultServer = "stun.l.google.com:19302"
)

type Client struct {
	serverAddr string
	StunName   string
	conn       net.PacketConn
}

type Host struct {
	family uint16
	port   uint16
	ip     string
}

type NAT int

func main() {
	client := NewClient()
	nat, host, err := client.NatHost()
	if err != nil {
		log.Print(err)
	}

	fmt.Println("NAT:", nat)
	if host != nil {
		fmt.Println("IP Family:", host.Family())
		fmt.Println("IP:", host.IP())
		fmt.Println("Port:", host.Port())
	}
}

func NewClient() *Client {
	client := new(Client)
	client.SetName(DefaultName)
	return client
}

func (c *Client) SetName(name string) {
	c.StunName = name
}

func (c *Client) SetServer(name string) {
	c.serverAddr = name
}

func (c *client) NatHost() (NAT, *Host, error) {
	c.SetServer(DefaultServer)
	serverUDPAddr, err := net.ResolveUDPAddr("udp", c.serverAddr)
	if err != nil {
		return nil, nil, err
	}

	co := c.conn
	co, err = net.ListenUDP("udp", nil)
	if err != nil {
		return nil, nil, err
	}
	defer co.Close()

}
