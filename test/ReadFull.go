package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"strings"
)

func main() {
	r := strings.NewReader("some io.Reader stream to be read\n")

	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf) // some

	if _, err := io.ReadFull(r, buf); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf) //  io.

	// minimal read size bigger than io.Reader stream
	// rのサイズがbufよりも少なければerr == EOFになる
	longBuf := make([]byte, 64) // len(r)
	if _, err := io.ReadFull(r, longBuf); err != nil {
		fmt.Println("error:", err) // error: unexpected EOF
	}

	var tbuf [12]byte
	if _, err := io.ReadFull(rand.Reader, tbuf[:]); err != nil { // make random byte
		log.Fatal(err)
	}
	fmt.Printf("%x\n", tbuf)
}
