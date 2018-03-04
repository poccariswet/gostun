package gostun

import "net"

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

func (addr *XORMappedAddr) GetAddr(m *Message, address AttributeType) error {
	return nil
}

func (addr *XORMappedAddr) GetXORMapped(m *Message) error {
	return addr.GetAddr(m, XOR_MAPPED_ADDRESS)
}
