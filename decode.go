package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	stunMessageHeader   = 20
	magicCookie         = 0x2112A442
	stunAttributeHeader = 4
)

type Client struct {
	conn Connection
}

type Message struct {
	Raw []byte
}

type Connection interface {
	io.Reader
}

func main() {
	conn, err := net.Dial("udp", "stun.l.google.com:19302")
	if err != nil {
		log.Fatal(err)
	}
	m := new(Message)
	c := &Client{
		conn: conn,
	}

	m.Raw = make([]byte, 1024)
	raw, err := m.ReadFrom(c.conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(raw)
}

func (m *Message) ReadFrom(r io.Reader) (int64, error) {
	tBuf := m.Raw[:cap(m.Raw)]

	n, err := r.Read(tBuf)
	if err != nil {
		return int64(n), err
	}
	m.Raw = tBuf[:n]
	return int64(n), m.Decode()
}

func (m *Message) Decode() error {
	buf := m.Raw
	if len(buf) < stunMessageHeader {
		return errors.New("unexpected EOF: not enough bytes to read header")
	}

	var (
		size        = int(bin.Uint16(buf[2:4]))
		magiccookie = bin.Uint32(buf[4:8])
		fullMsg     = messageHeaderSize + size
	)

	if magiccookie != magicCookie {
		err := fmt.Sprintf("%x is invalid magic cookie. maegic Cookie: %x", magiccookie, magicCookie)
		return errors.New(err)
	}

	if len(buf) < fullMsg {
		err := fmt.Sprintf("buffer length is invalid %d, expected message size %d", len(buf), fullMsg)
		return errors.New(err)
	}

	//copy header data

	return nil
}
