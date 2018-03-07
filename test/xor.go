package main

import "fmt"

func main() {
	a := 5 // 101
	b := 6 // 110

	fmt.Printf("%b\n", a)
	fmt.Printf("%b\n", b)
	fmt.Printf("%b\n", a^b) //11
}
