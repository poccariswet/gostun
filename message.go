package gostun

import (
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
