package pool

import (
	"testing"

	"github.com/peienxie/socketpool/echo"
	"github.com/stretchr/testify/require"
)

func TestSocketPool(t *testing.T) {
	// Start local echo server for testing
	listener, err := echo.StartEchoServer("127.0.0.1:0")
	require.NoError(t, err, "Failed to start echo server")
	defer listener.Close()

	pool, err := NewSocketPool(2, listener.Addr().String())
	require.NoError(t, err)
	defer pool.Close()

	conn1, err := pool.Get()
	require.NoError(t, err)
	defer pool.Release(conn1)
	// Check that conn1 is connected to the correct address
	require.Equal(t, listener.Addr().String(), conn1.RemoteAddr().String(),
		"Unexpected connection address: want %s but got %s", listener.Addr().String(), conn1.RemoteAddr().String())

	conn2, err := pool.Get()
	require.NoError(t, err)
	defer pool.Release(conn2)
	// Check that conn2 is not the same as conn1
	require.NotEqual(t, conn1, conn2, "Expected different connections, got the same connection")

	conn3, err := pool.Get()
	require.Error(t, err, "Expected error when getting connection from full pool, got none")
	require.Nil(t, conn3)

	// Test that connection pool size is correct
	require.Equalf(t, 2, pool.Active(), "Unexpected pool size: want 2 but got %d", pool.Active())
}
