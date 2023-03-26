package factory

import (
	"crypto/tls"
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

type TLSSocketFactory struct {
	addr   string
	config *tls.Config
}

func NewTLSSocketFactory(addr string, config *tls.Config) TLSSocketFactory {
	return TLSSocketFactory{addr: addr, config: config}
}

func (factory TLSSocketFactory) CreateSocket() (net.Conn, error) {
	return tls.Dial("tcp", factory.addr, factory.config)
}
