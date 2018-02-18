package gostun

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	magicCookie       = 0x2112A442
	TransactionIDSize = 12 // 96 bit
	messageHeader     = 20
	attributeHeader   = 4
)

type MessageClass byte

type Method uint16

const (
	Request         MessageClass = 0x00 // 0b00
	Indication      MessageClass = 0x01 // 0b01
	SuccessResponse MessageClass = 0x02 // 0b10
	ErrorResponse   MessageClass = 0x03 // 0b11
)

// MessageType is STUN Message Type Field.
type MessageType struct {
	Method Method       // binding
	Class  MessageClass // request
}

type Message struct {
	Raw           []byte
	Type          MessageType
	Length        uint32
	TransactionID [TransactionIDSize]byte
}

func (m *Message) ReadConn(r io.Reader) (int, error) {
	raw := m.Raw

	n, err := r.Read(raw)
	if err != nil {
		return n, err
	}

	m.Raw = raw[:n]
	//	return n, m.Decode()
	return n, nil
}

func (m *Message) Decode() error {
	header := m.Raw
	mtype := binary.BigEndian.Uint16(header[0:2])   //STUN Message type
	mlength := binary.BigEndian.Uint16(header[2:4]) //STUN Message length
	mcookie := binary.BigEndian.Uint16(header[4:8]) //Magic Cookie
	fullHeader := messageHeaderSize + int(mlength)  //len(m.Raw)

	if mcookie != magicCookie {
		err := fmt.Sprintf("%x is invalid value magicCookie is %x\n", mcookie, magicCookie)
		return errors.New(err)
	}

	if len(header) < fullHeader {
		err := fmt.Sprintf("this length %d is less than %d", len(header), fullHeader)
		return errors.New(err)
	}

	m.Type.ReadValue(mtype)    // copy STUN message type
	m.Length = uint32(mlength) // copy STUN message type

	return nil
}

const (
	bic0    = 0x1
	bic1    = 0x2
	shiftc0 = 4
	shiftc1 = 7
)

/*
    0                 1
    2  3  4 5 6 7 8 9 0 1 2 3 4 5

   +--+--+-+-+-+-+-+-+-+-+-+-+-+-+
   |M |M |M|M|M|C|M|M|M|C|M|M|M|M|
   |11|10|9|8|7|1|6|5|4|0|3|2|1|0|
   +--+--+-+-+-+-+-+-+-+-+-+-+-+-+
									7					4
   Format of STUN Message Type Field


*/

func (mt *MessageType) ReadValue(v uint16) {
	// difine class
	c0 := (v >> shiftc0) & bitc0
	c1 := (v >> shiftc1) & bitc1
	Class := c0 + c1
	mt.Class = MessageClass(Class)

}
