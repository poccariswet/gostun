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

const (
	Request         MessageClass = 0x00 // 0b00
	Indication      MessageClass = 0x01 // 0b01
	SuccessResponse MessageClass = 0x02 // 0b10
	ErrorResponse   MessageClass = 0x03 // 0b11
)

// STUN Message Type Field.
type MessageType struct {
	Method Method       // binding
	Class  MessageClass // request
}

type Message struct {
	Raw           []byte
	Type          MessageType
	Length        uint32
	TransactionID [TransactionIDSize]byte
	Attributes    Attributes
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

/*
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |         Type                  |            Length             |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                         Value (variable)                ....
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

                         Format of STUN Attributes
*/

type AttributeType uint16

type AttributeField struct {
	Type   AttributeType
	Length uint16 // ignored while encoding
	Value  []byte
}

type Attributes []AttributeField

// Comprehension-required range (0x0000-0x7FFF): page 43
const (
	MAPPED_ADDRESS     AttributeType = 0x0001
	USERNAME           AttributeType = 0x0006
	MESSAGE_INTEGRITY  AttributeType = 0x0008
	ERROR_CODE         AttributeType = 0x0009
	UNKNOWN_ATTRIBUTES AttributeType = 0x000A
	REALM              AttributeType = 0x0014
	NONCE              AttributeType = 0x0015
	XOR_MAPPED_ADDRESS AttributeType = 0x0020

	SOFTWARE         AttributeType = 0x8022
	ALTERNATE_SERVER AttributeType = 0x8023
	FINGERPRINT      AttributeType = 0x8028
)

func (m *Message) Decode() error {
	header := m.Raw
	mtype := binary.BigEndian.Uint16(header[0:2])   //STUN Message type
	mlength := binary.BigEndian.Uint16(header[2:4]) //STUN Message length
	mcookie := binary.BigEndian.Uint32(header[4:8]) //Magic Cookie
	fullHeader := messageHeader + int(mlength)      //len(m.Raw)

	if mcookie != magicCookie {
		err := fmt.Sprintf("%x is invalid value magicCookie is %x\n", mcookie, magicCookie)
		return errors.New(err)
	}

	if len(header) < fullHeader {
		err := fmt.Sprintf("this length %d is less than %d", len(header), fullHeader)
		return errors.New(err)
	}

	m.Type.ReadValue(mtype)                           // copy STUN message type
	m.Length = uint32(mlength)                        // copy STUN message type
	copy(m.TransactionID[:], header[8:messageHeader]) // copy STUN Transaction ID (96 bits|12 byte)

	m.Attributes = m.Attributes[:0]
	buf := header[messageHeader:fullHeader]
	attrsize := 0 // initialize

	for attrsize < int(mlength) {
		// github.com/soeyusuke/note/stun => STUN Attributes
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
		buf = buf[attributeHeader:] // adjust buf to Value
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

func (mt *MessageType) ReadValue(v uint16) {
	// difine class
	c0 := (v >> shiftc0) & bitc0
	c1 := (v >> shiftc1) & bitc1
	Class := c0 + c1
	mt.Class = MessageClass(Class)

	// method
	m0m3 := v & mbit1
	m4m6 := (v >> methodshift1) & mbit2
	m7m11 := (v >> methodshift2) & mbit3

	m := m0m3 + m4m6 + m7m11
	mt.Method = Method(m)
}

// Since STUN aligns attributes on 32-bit boundaries, attributes whose content
// is not a multiple of 4 bytes are padded with 1, 2, or 3 bytes of
// padding so that its value contains a multiple of 4 bytes.  The
// padding bits are ignored, and may be any value.
func (a *AttributeField) PaddingValue() int {
	const padding = 4
	al := int(a.Length)
	l := padding * (al / padding)
	if l < al {
		l += padding
	}
	return l
}
