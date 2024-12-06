package sorter

import "net"

type Packet struct {
	SourceIP   net.IP
	DestIP     net.IP
	SourcePort int
	DestPort   int
	Result     string
	Protocol   string
}
