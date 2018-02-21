package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	b := []byte{0xe8, 0x03, 0xd0, 0x07}
	x1 := binary.BigEndian.Uint16(b[0:])
	x2 := binary.BigEndian.Uint16(b[2:])
	fmt.Println("BigEndian")
	fmt.Printf("%#04x %#04x\n", x1, x2)

	x3 := binary.LittleEndian.Uint16(b[0:])
	x4 := binary.LittleEndian.Uint16(b[2:])
	fmt.Println("LittleEndian")
	fmt.Printf("%#04x %#04x\n", x3, x4)
}

//go run bin.go
//BigEndian
//0xe803 0xd007
//LittleEndian
//0x03e8 0x07d0
