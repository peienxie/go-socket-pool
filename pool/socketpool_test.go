package pool

import (
	"crypto/tls"
	"testing"

	"github.com/peienxie/socketpool/echo"
	"github.com/peienxie/socketpool/factory"
	"github.com/stretchr/testify/require"
)

func TestSocketPool(t *testing.T) {
	// Start local echo server for testing
	listener, err := echo.StartEchoServer("127.0.0.1:0")
	require.NoError(t, err, "Failed to start echo server")
	defer listener.Close()

	factory := factory.NewSocketFactory(listener.Addr().String())
	socketPool, err := NewSocketPool(100, factory)
	require.NoError(t, err)
	defer socketPool.Close()
	testSocketPool(t, socketPool, listener.Addr().String())
}

func TestTLSSocketPool(t *testing.T) {
	// Start local echo server for testing
	cert, err := tls.LoadX509KeyPair("../certs/server.crt", "../certs/server.key")
	require.NoError(t, err)
	serverConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := echo.StartTLSEchoServer("127.0.0.1:0", serverConfig)
	require.NoError(t, err, "Failed to start echo server")
	defer listener.Close()

	config := &tls.Config{InsecureSkipVerify: true}
	factory := factory.NewTLSSocketFactory(listener.Addr().String(), config)
	socketPool, err := NewSocketPool(100, factory)
	require.NoError(t, err)
	defer socketPool.Close()
	testSocketPool(t, socketPool, listener.Addr().String())
}

func testSocketPool(t *testing.T, socketPool *SocketPool, addr string) {
	conn1, err := socketPool.Get()
	require.NoError(t, err)
	defer socketPool.Release(conn1)
	// Check that conn1 is connected to the correct address
	require.Equal(t, addr, conn1.RemoteAddr().String(),
		"Unexpected connection address: want %s but got %s", addr, conn1.RemoteAddr().String())

	for i := 0; i < socketPool.Size()-1; i++ {
		conn2, err := socketPool.Get()
		require.NoError(t, err)
		defer socketPool.Release(conn2)
		// Check that conn2 is not the same as conn1
		require.NotEqual(t, conn1, conn2, "Expected different connections, got the same connection")
	}

	conn3, err := socketPool.Get()
	require.Error(t, err, "Expected error when getting connection from full pool, got none")
	require.Nil(t, conn3)

	// Test that connection pool size is correct
	require.Equalf(t, socketPool.Size(), socketPool.Active(), "Unexpected pool active count: want %d but got %d", socketPool.Size(), socketPool.Active())
}
