package gostun

import "io"

const (
	magicCookie       = 0x2112A442
	TransactionIDSize = 12 // 96 bit
	messageHeader     = 20
	attributeHeader   = 4
)

type MessageClass byte

const (
	Request         MessageClass = 0x00 // 0b00
	Indication      MessageClass = 0x01 // 0b01
	SuccessResponse MessageClass = 0x02 // 0b10
	ErrorResponse   MessageClass = 0x03 // 0b11
)

type Message struct {
	Type          MessageType
	Length        uint32
	TransactionID [TransactionIDSize]byte
	Attributes    Attributes
	Raw           []byte
}

func (m *Message) ReadConn(r io.Reader) (int, error) {
	raw := m.Raw

	n, err := r.Read(raw)
	if err != nil {
		return n, err
	}

	m.Raw = raw[:n]
	return n, m.Decode()
}
