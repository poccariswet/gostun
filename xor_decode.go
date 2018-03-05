package gostun

import (
	"errors"
	"net"
)

type Addr struct {
	Port int    // port number
	IP   net.IP // IP address
}

/*
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |x x x x x x x x|    Family     |         X-Port                |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                X-Address (Variable)
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
               Format of XOR-MAPPED-ADDRESS Attribute
*/

type XORMappedAddr Addr

func (m *Message) GetAttrFiledValue(attrtype AttributeType) ([]byte, error) {
	for _, attr := range m.Attributes {
		if attr.Type == attrtype {
			return attr.Value, nil
		}
	}
	return nil, errors.New("Attribute is not matched")
}

func (addr *XORMappedAddr) DecodexorAddr(m *Message, attrtype AttributeType) error {
	val, err := m.GetAttrFiledValue(attrtype)
	if err != nil {
		return err
	}

	return nil
}

func (addr *XORMappedAddr) GetXORMapped(m *Message) error {
	return addr.DecodexorAddr(m, XOR_MAPPED_ADDRESS)
}
