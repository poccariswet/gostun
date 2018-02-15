package main

import (
	"fmt"
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
	fmt.Println(c)
}
