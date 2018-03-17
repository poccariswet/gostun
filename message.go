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
	attributeHeader   = 4 // type and length
)

type MessageClass byte

type Method uint16

// class
const (
	Request         MessageClass = 0x00 // 0b00
	Indication      MessageClass = 0x01 // 0b01
	SuccessResponse MessageClass = 0x02 // 0b10
	ErrorResponse   MessageClass = 0x03 // 0b11
)

// method
const (
	MethodBinding Method = 0x001
)

// Binding Message type
var (
	BindingRequest = NewMessageType(MethodBinding, Request)
	BindingSuccess = NewMessageType(MethodBinding, SuccessResponse)
	BindingError   = NewMessageType(MethodBinding, ErrorResponse)
)

// STUN Message Type Field.
type MessageType struct {
	Method Method       // binding
	Class  MessageClass // request
}

//reutn new message type has Method and Class
func NewMessageType(m Method, c MessageClass) MessageType {
	return MessageType{
		Method: m,
		Class:  c,
	}
}

// stun message type
type Message struct {
	Raw           []byte //full message
	Type          MessageType
	Length        uint32
	TransactionID [TransactionIDSize]byte
	Attributes    Attributes
}

func (m *Message) ReadConn(r io.Reader) (int, error) {
	n, err := r.Read(m.Raw)
	if err != nil {
		return n, err
	}
	m.Raw = m.Raw[:n]

	return n, m.Decode()
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

func (m *Message) Decode() error {
	header := m.Raw
	mtype := binary.BigEndian.Uint16(header[0:2])   //STUN Message type
	mlength := binary.BigEndian.Uint16(header[2:4]) //STUN Message length
	mcookie := binary.BigEndian.Uint32(header[4:8]) //Magic Cookie
	fullHeader := messageHeader + int(mlength)      //len(m.Raw)

	// check magic cookie
	if mcookie != magicCookie {
		err := fmt.Sprintf("%x is invalid value magicCookie is %x\n", mcookie, magicCookie)
		return errors.New(err)
	}
	// check header size
	if len(header) < fullHeader {
		err := fmt.Sprintf("this length %d is less than %d", len(header), fullHeader)
		return errors.New(err)
	}

	m.Type.DecodeMessageType(mtype)                   // copy STUN message type
	m.Length = uint32(mlength)                        // copy STUN message type
	copy(m.TransactionID[:], header[8:messageHeader]) // copy STUN Transaction ID (96 bits|12 byte)
	if err := m.AttrDecode(header[messageHeader:fullHeader], int(mlength)); err != nil {
		return err
	}

	return nil
}

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

/*
    0                 1
    2  3  4 5 6 7 8 9 0 1 2 3 4 5

   +--+--+-+-+-+-+-+-+-+-+-+-+-+-+
   |M |M |M|M|M|C|M|M|M|C|M|M|M|M|
   |11|10|9|8|7|1|6|5|4|0|3|2|1|0|
   +--+--+-+-+-+-+-+-+-+-+-+-+-+-+
                  7       4
   Format of STUN Message Type Field
*/

func (mt *MessageType) DecodeMessageType(v uint16) {
	// difine class
	c0 := (v >> shiftc0) & bitc0
	c1 := (v >> shiftc1) & bitc1
	class := c0 + c1
	mt.Class = MessageClass(class)

	// method
	m0m3 := v & mbit1
	m4m6 := (v >> methodshift1) & mbit2
	m7m11 := (v >> methodshift2) & mbit3

	m := m0m3 + m4m6 + m7m11
	mt.Method = Method(m)
}

func (m *Message) AttrDecode(buf []byte, l int) error {
	m.Attributes = m.Attributes[:0]
	attrsize := 0 // initialize

	for attrsize < l {
		attr := AttributeField{
			Type:   AttributeType(binary.BigEndian.Uint16(buf[0:2])), //Attribute type - first 2byte
			Length: binary.BigEndian.Uint16(buf[2:4]),                // Attributes Length - next 2byte
		}

		if len(buf) < attributeHeader {
			err := fmt.Sprintf("buf(%d) is less than attributeHeader(%d)", len(buf), attributeHeader)
			return errors.New(err)
		}

		alen := attr.PaddingValue() // padding
		attrsize += attributeHeader // increment attrsize 4byte(type + length)
		buf = buf[attributeHeader:] // adjust 4 byte buf to Value
		if len(buf) < alen {
			err := fmt.Sprintf("buf length(%d) is less than value size is expected(%d)", len(buf), alen)
			return errors.New(err)
		}

		attr.Value = buf[:int(attr.Length)]
		attrsize += alen // increment attrsize Value size
		buf = buf[alen:] // adjust buf Attribute Field

		m.Attributes = append(m.Attributes, attr)
	}
	return nil
}
