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
	ip string
}

type NAT int

func main() {
	client := NewClient()
	Nat, Host, err := client.GetNatHost()
	if err != nil {
		log.Print(err)
	}

	fmt.Println("NAT:", Nat)
	if host != nil {
		fmt.Println("External IP: %v", Host.IP())
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

func (c *Client) GetNatHost() (NAT, *Host, error) {
	c.SetServer(DefaultServer)
	addr, err := net.ResolveUDPAddr("udp", c.serverAddr)
	if err != nil {
		return nil, nil, err
	}

	co := c.conn
	if co != nil {
		co, err = net.ListenUDP("udp", nil)
		if err != nil {
			return nil, nil, err
		}
	}
	defer co.Close()

	resp, err := c.SendReq(conn, addr)
	if err != nil {
		return nil, nil, err
	}

}
