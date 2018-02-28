package gostun

import (
	"encoding/binary"
	"log"
)

type transactionIDSetter struct{}

func (transactionIDSetter) AddTo(m *Message) error {
	return nil
}

// Sets *Message attribute.
type MsgSetter interface {
	AddTo(m *Message) error
}

var TransactionID MsgSetter = transactionIDSetter{}

// reset message
func (m *Message) Reset() {
	m.Raw = m.Raw[:0]
	m.Length = 0
	m.Attributes = m.Attributes[:0]
}

func (m *Message) AllocRaw() {
	l := len(m.Raw) + messageHeader
	for cap(m.Raw) < l {
		m.Raw = append(m.Raw, 0)
	}
	m.Raw = m.Raw[:l]
}

/*
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |0 0|     STUN Message Type     |         Message Length        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                         Magic Cookie                          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                                                               |
   |                     Transaction ID (96 bits)                  |
   |                                                               |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	                   Format of STUN Message Header
*/

/*
    0                 1
    2  3  4 5 6 7 8 9 0 1 2 3 4 5

   +--+--+-+-+-+-+-+-+-+-+-+-+-+-+
   |M |M |M|M|M|C|M|M|M|C|M|M|M|M|
   |11|10|9|8|7|1|6|5|4|0|3|2|1|0|
   +--+--+-+-+-+-+-+-+-+-+-+-+-+-+
                  7       4
   Format of STUN Message Type Field

	const (
	bitc0   = 0x1
	bitc1   = 0x2
	shiftc0 = 4
	shiftc1 = 7

	methodshift1 = 1
	methodshift2 = 2
	mbit1        = 0xf   //M0~M3=>0b0000000000001111
	mbit2        = 0x70  //M4~M6=>0b0000000001110000
	mbit3        = 0xf80 //M7~M11=>0b00011111000000
)

*/

// write message type to m.Raw
func (m *Message) WriteMessageType() {
	// Class
	class := uint16(m.Type.Class)
	c0 := (class & bitc0) << shiftc0 // 4 bit shift
	c1 := (class & bitc1) << shiftc1 // 7 bit shift
	c := c0 + c1

	// Method
	method := uint16(m.Type.Method)
	m1m3 := method & mbit1
	m4m6 := method & mbit2
	m7m11 := method & mbit3
	method = m1m3 + (m4m6 << methodshift1) + (m7m11 << methodshift2)

	mtype := c + method

	binary.BigEndian.PutUint16(m.Raw[0:2], mtype)
}

func (m *Message) WriteMessageLength() {
	binary.BigEndian.PutUint16(m.Raw[2:4], m.Length)
}

func (m *Message) WriteMagicCookie() {
	binary.BigEndian.PutUint16(m.Raw[4:8], magicCookie)
}

func (m *Message) WriteTransactionID() {
	copy(m.Raw[8:messageHeader], m.TransactionID[:])
}

// make message header
func (m *Message) WriteMessageHeader() {
	m.AllocRaw() // alloc 0, part of message header size
	m.WriteMessageType()
	m.WriteMessageLength()
	m.WriteMagicCookie()
	m.WriteTransactionID()

}

func (m *Message) build(s ...MsgSetter) error {
	m.Reset()
	m.WriteMessageHeader()

	return nil
}

// wraps m.build
func Build(s ...MsgSetter) (*Message, error) {
	m := new(Message)
	return m, m.build(s...)
}

func MessageBuild(s ...MsgSetter) *Message {
	m, err := Build(s...)
	if err != nil {
		log.Fatal(err)
	}

	return m
}
