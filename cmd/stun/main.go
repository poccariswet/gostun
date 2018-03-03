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
	_ = c

	m := gostun.MessageBuild(gostun.TransactionID, gostun.BindingRequest)
	fmt.Printf("Type: %x, %x\n", m.Type.Class, m.Type.Method)
	fmt.Printf("Length: %v\n", m.Length)
	fmt.Printf("Magic cookie: %x\n", m.Raw[4:8])
	fmt.Printf("TransactionID: %x\n", m.TransactionID)
	fmt.Printf("Attributes: %v\n", m.Attributes)
	fmt.Printf("Raw: %v\n", m.Raw)

	f := func(event gostun.EventObject) {
		if event.Err != nil {
			log.Fatal(err)
		}
	}

	rto := time.Now().Add(time.Second * 5)

	if err := c.Call(m, f, rto); err != nil {
		log.Fatal(err)
	}
}
