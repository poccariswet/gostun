package main

import (
	"fmt"
	"log"
	"time"

	"github.com/soeyusuke/gostun"
)

const (
	defaultAddr = "stun.l.google.com:19302"
)

func main() {
	c, err := gostun.Dial("udp", defaultAddr)
	if err != nil {
		log.Fatal(err)
	}

	m := gostun.MessageBuild(gostun.TransactionID, gostun.BindingRequest)

	f := func(msg gostun.MessageObj) {
		var addr gostun.XORMappedAddr
		if msg.Err != nil {
			log.Fatal(err)
		}
		if err := addr.GetXORMapped(msg.Msg); err != nil {
			log.Fatal(err)
		}
		fmt.Println(addr)
	}

	rto := time.Now().Add(time.Second * 5)

	if err := c.Call(m, f, rto); err != nil {
		log.Fatal(err)
	}
}
