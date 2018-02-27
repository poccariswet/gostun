package main

import (
	"fmt"
)

const (
	messageHeaderSize = 20
)

type Message struct {
	Raw []byte
	Ex  []int
}

func (m *Message) AllocRaw(v int) {
	n := len(m.Raw) + v
	for cap(m.Raw) < n {
		m.Raw = append(m.Raw, 0)
		m.Ex = append(m.Ex, 0)
	}
	m.Raw = m.Raw[:n]
	m.Ex = m.Ex[:n]
	copy(m.Ex, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
}

func main() {
	m := new(Message)
	m.Raw = m.Raw[:0]

	if len(m.Raw) < messageHeaderSize {
		fmt.Println(len(m.Raw)) // 0
		fmt.Println(cap(m.Raw)) // 0
		m.AllocRaw(messageHeaderSize)
	}
	fmt.Println(m.Ex)             // [1 2 3 4 5 6 7 8 9 10 0 0 0 0 0 0 0 0 0 0]
	m.Ex = m.Ex[:3]               //
	fmt.Println(len(m.Ex))        // 3
	fmt.Println(cap(m.Ex))        // 32
	fmt.Println(m.Ex)             // [1 2 3]
	fmt.Println(m.Raw)            // [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	_ = m.Raw[:messageHeaderSize] //messageHeaderSize分割り当てがないと、パニック
	fmt.Println(m.Raw)            // [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	fmt.Println(len(m.Raw))       // 20
	fmt.Println(cap(m.Raw))       // 32
}
