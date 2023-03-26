package factory

import (
	"net"
)

type SocketFactory struct {
	addr string
}

func NewSocketFactory(addr string) SocketFactory {
	return SocketFactory{addr: addr}
}

func (factory SocketFactory) CreateSocket() (net.Conn, error) {
	return net.Dial("tcp", factory.addr)
}
