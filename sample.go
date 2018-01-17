package main

import (
	"fmt"
	"net"
	"os"
)

const (
	DefaultAddr   = "www.google.com"
	DefaultServer = "stun.l.google.com:19302"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", DefaultServer)
	if err != nil {
		fmt.Println("Resolve error ")
		os.Exit(1)
	}
	fmt.Println("Addr is ", addr.String())
	fmt.Println("IP ", addr.IP)
	fmt.Println("Port ", addr.Port)
	fmt.Println("Zone ", addr.Zone)
}
