package pool

import (
	"io"
	"net"
	"testing"

	"github.com/peienxie/socketpool/echo"
)

const SERVER_ADDR = "127.0.0.1:0"

func BenchmarkNormalTcpConnectionWithSocketPool(b *testing.B) {
	// Start local echo server for testing
	listener, err := echo.StartEchoServer(SERVER_ADDR)
	if err != nil {
		b.Fatalf("Failed to start echo server: %v", err)
	}
	defer listener.Close()

	// Create socket pool
	pool, err := NewSocketPool(8, listener.Addr().String())
	if err != nil {
		b.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	// Prepare data for benchmark
	data := createData()

	// Run benchmark
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, err := pool.Get()
			if err != nil {
				b.Fatalf("Failed to get connection from pool: %v", err)
			}
			if _, err := conn.Write(data); err != nil {
				b.Fatalf("Failed to write data: %v", err)
			}
			if _, err := io.ReadFull(conn, data); err != nil {
				b.Fatalf("Failed to read data: %v", err)
			}
			pool.Release(conn)
		}
	})
}

func BenchmarkNormalTcpConnectionWithoutSocketPool(b *testing.B) {
	listener, err := echo.StartEchoServer(SERVER_ADDR)
	if err != nil {
		b.Fatalf("Failed to start echo server: %v", err)
	}
	defer listener.Close()

	// Prepare data for benchmark
	data := createData()

	// Run benchmark
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Connect to server
			conn, err := net.Dial("tcp", listener.Addr().String())
			if err != nil {
				b.Fatalf("Failed to dial: %v", err)
			}

			// Write data
			if _, err := conn.Write(data); err != nil {
				b.Fatalf("Failed to write data: %v", err)
			}

			// Read echoed data
			if _, err := io.ReadFull(conn, data); err != nil {
				b.Fatalf("Failed to read data: %v", err)
			}

			// Close connection
			if err := conn.Close(); err != nil {
				b.Fatalf("Failed to close connection: %v", err)
			}
		}
	})
}

func createData() []byte {
	data := make([]byte, 1024)
	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}
	return data
}
