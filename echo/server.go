package echo

import (
	"crypto/tls"
	"io"
	"net"
)

func StartEchoServer(addr string) (net.Listener, error) {
	// Create listener to listen on given address
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Start echo server
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c)
			}(conn)
		}
	}()

	// Return server listener
	return listener, nil
}

func StartTLSEchoServer(addr string, config *tls.Config) (net.Listener, error) {
	// Create a listener to listen on given address
	listener, err := tls.Listen("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	// Start echo server
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			go func(c net.Conn) {
				defer conn.Close()
				io.Copy(c, c)
			}(conn)
		}
	}()

	// Return server listener
	return listener, nil
}
