package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	l := make([]byte, 4)
	binary.LittleEndian.PutUint16(l[0:], 0x03e8) // put the l value
	binary.LittleEndian.PutUint16(l[2:], 0x07d0)
	fmt.Printf("% x\n", l)

	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[0:], 0x03e8)
	binary.BigEndian.PutUint16(b[2:], 0x07d0)
	fmt.Printf("% x\n", b)
}
