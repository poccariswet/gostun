package main

import (
	"fmt"
	"net"
)

type IPV interface {
	Ensure(ipLen int)
}

type Addr struct {
	Port int
	IP   net.IP
}

type Ipv4 Addr
type Ipv6 Addr

func (a *Ipv4) String() string {
	str := fmt.Sprintf("IPv4\nIP: %v\nlen(IP): %d\ncap(IP): %d", a.IP, len(a.IP), cap(a.IP))
	return str
}

func (a *Ipv6) String() string {
	str := fmt.Sprintf("IPv6\nIP: %v\nlen(IP): %d\ncap(IP): %d", a.IP, len(a.IP), cap(a.IP))
	return str
}

func (a *Ipv4) Ensure(ipLen int) {
	ipLen = net.IPv4len
	if len(a.IP) < ipLen {
		a.IP = a.IP[:cap(a.IP)]
		for len(a.IP) < ipLen {
			a.IP = append(a.IP, 0)
		}
	}
}

func (a *Ipv6) Ensure(ipLen int) {
	ipLen = net.IPv6len
	if len(a.IP) < ipLen {
		a.IP = a.IP[:cap(a.IP)]
		for len(a.IP) < ipLen {
			a.IP = append(a.IP, 0)
		}
	}
}

func main() {
	var ipv4 IPV = new(Ipv4)
	var ipv6 IPV = new(Ipv6)

	ipv4.Ensure(0)
	ipv6.Ensure(0)
	fmt.Println(ipv4)
	fmt.Println(ipv6)
}
