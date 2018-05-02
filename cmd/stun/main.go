package main

import (
	"fmt"
	"log"
	"time"

	"github.com/soeyusuke/gostun"
)

const (
	defaultAddr = "stun.l.google.com:19302" //ICE server
)

func main() {
	c, err := gostun.Dial("udp", defaultAddr)
	if err != nil {
		log.Fatal(err)
	}

	m := gostun.MessageBuild(gostun.TransactionID, gostun.BindingRequest)
	rto := time.Now().Add(time.Second * 5)

	addr, err := c.Call(m, rto)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(addr)
}
