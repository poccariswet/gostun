package main

import (
	"log"
)

type Client struct {
	StunName string
}

func main() {
	client := NewClient()
	nat, host, err := client.NatHost()
	if err != nil {
		log.Print(err)
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

func (c *client) NatHost() {

}
