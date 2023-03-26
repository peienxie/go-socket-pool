package pool

import (
	"fmt"
	"net"
	"sync"
)

type SocketPool struct {
	mu     sync.Mutex
	conns  []net.Conn
	size   int
	factory SocketFactory
	active int
}

type SocketFactory interface {
	CreateSocket() (net.Conn, error)
}

func NewSocketPool(size int, factory SocketFactory) (*SocketPool, error) {
	conns := make([]net.Conn, size)
	return &SocketPool{
		conns:  conns,
		size:   size,
		active: 0,
		factory: factory,
	}, nil
}

func (p *SocketPool) Size() int {
	return p.size
}

func (p *SocketPool) Active() int {
	return p.active
}

func (p *SocketPool) Get() (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active >= p.size {
		return nil, fmt.Errorf("no available connections in the pool")
	}
	conn := p.conns[p.active]
	if conn == nil {
		var err error
		conn, err = p.factory.CreateSocket()
		if err != nil {
			return nil, err
		}
		p.conns[p.active] = conn
	}
	p.active++
	return conn, nil
}

func (p *SocketPool) Release(conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.active--
}

func (p *SocketPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, conn := range p.conns {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
