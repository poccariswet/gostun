package main

import (
	"log"

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

	//	rto := time.Now().Add(time.Second * 5)
	if err := gostun.MessageBuild(gostun.TransactionID, gostun.BindingRequest); err != nil {
		log.Fatal(err)
	}
}
