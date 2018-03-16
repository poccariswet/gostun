package gostun

import (
	"encoding/binary"
	"errors"
	"fmt"
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

//  0x01:IPv4
//  0x02:IPv6
const (
	IPv4 uint16 = 0x01
	IPv6 uint16 = 0x02
)

type XORMappedAddr Addr

func (addr *XORMappedAddr) String() string {
	return fmt.Sprintf("IP: %v\nPort:%d\n", addr.IP, addr.Port)
}

func (attr Attributes) GetAttrFiledValue(attrtype AttributeType) ([]byte, error) {
	for _, a := range attr {
		fmt.Println(a.Type, "\n")
		if a.Type == attrtype {
			return a.Value, nil
		}
	}
	return nil, errors.New("Attribute is not matched")
}

func (m *Message) GetRapped(attrtype AttributeType) ([]byte, error) {
	val, err := m.Attributes.GetAttrFiledValue(attrtype)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (addr *XORMappedAddr) DecodexorAddr(m *Message, attrtype AttributeType) error {
	val, err := m.GetRapped(attrtype)
	if err != nil {
		return err
	}

	var (
		family uint16
		ipl    int
	)
	family = binary.BigEndian.Uint16(val[0:2])

	//check family address
	if family == IPv4 {
		ipl = net.IPv4len
	} else if family == IPv6 {
		ipl = net.IPv6len
	} else {
		err := fmt.Sprintf("family decode err: family = %d\n", family)
		return errors.New(err)
	}

	addr.IP = addr.IP[:cap(addr.IP)]
	for len(addr.IP) < ipl {
		addr.IP = append(addr.IP, 0)
	}
	addr.IP = addr.IP[:ipl]

	/*
		X-Port is computed by taking the mapped port in host byte order,
		 XOR'ing it with the most significant 16 bits of the magic cookie, and
		 then the converting the result to network byte order.
	*/

	/*
		If the IPaddress family is IPv4, X-Address is computed by taking the mapped IP
		address in host byte order, XOR'ing it with the magic cookie, and
		converting the result to network byte order.  If the IP address
		family is IPv6, X-Address is computed by taking the mapped IP address
		in host byte order, XOR'ing it with the concatenation of the magic
		cookie and the 96-bit transaction ID, and converting the result to
		network byte order(BigEndian).

		IPアドレスファミリーがIPv4であるなら、X-アドレスは、マップされたIPアドレスをホ
		ストバイトオーダーで取得して、マジッククッキーでそれをXORして、そして
		その結果をネットワークバイトオーダーに変換することで計算される

		IPアドレスファミリーがIPv6であるなら、X-アドレスは、マップされたIPアドレ
		スをホストバイトオーダーで取得して、マジッククッキーと96ビットのトラ
		ンザクションIDとを連結したものでそれをXORして、そしてその結果をネット
		ワークバイトオーダーに変換することで計算される
	*/
	addr.XorAddr(m.TransactionID[:], val)

	return nil
}

func (addr *XORMappedAddr) XorAddr(trans, value []byte) {
	//	xaddress := make([]byte, 4+TransactionIDSize)

	for i := 0; i < len(value)-4; i++ {
		addr.IP = append(addr.IP, value[i+4]^trans[i])
	}

	mscookie := magicCookie >> 16
	addr.Port = int(binary.BigEndian.Uint16(value[2:4])) ^ mscookie
}

func (addr *XORMappedAddr) GetXORMapped(m *Message) error {
	return addr.DecodexorAddr(m, XOR_MAPPED_ADDRESS)
}
