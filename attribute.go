package gostun

import "fmt"

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

var AttrTypeName = map[AttributeType]string{
	MAPPED_ADDRESS:     "MAPPED-ADDRESS",
	USERNAME:           "USERNAME",
	MESSAGE_INTEGRITY:  "MESSAGE-INTEGRITY",
	ERROR_CODE:         "ERROR-CODE",
	UNKNOWN_ATTRIBUTES: "UNKNOWN-ATTRIBUTES",
	REALM:              "REALM",
	NONCE:              "NONCE",
	XOR_MAPPED_ADDRESS: "XOR-MAPPED-ADDRESS",

	SOFTWARE:         "SOFTWARE",
	ALTERNATE_SERVER: "ALTERNATE_SERVER",
	FINGERPRINT:      "FINGERPRINT",
}

func (at AttributeType) String() string {
	name, ok := AttrTypeName[at]
	if !ok {
		return fmt.Sprintf("non-attribute: %x", at)
	}
	return name
}

func (af AttributeField) String() string {
	return fmt.Sprintf("%s: 0x%x", af.Type, af.Value)
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
