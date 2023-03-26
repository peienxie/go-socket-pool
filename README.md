# Socket Pool

This is a simple implementation of a socket pool in Go, which can be used to manage a pool of network socket connections for efficient reuse.

## Features

- Supports both TCP and TLS connections
- Lazy connection creation: connections are only created when needed
- Configurable connection pool capacity
- Thread-safe: can be used concurrently by multiple goroutines

## Installation

To use this package in your Go project, you can simply run:

```sh
go get github.com/peienxie/go-socket-pool
```

## Usage

To use the socket pool, you need to create a factory that implements `SocketFactory` interface, which is used to create new network socket connections.
And use this factory to create a socket pool with a specific size.

```
type SocketFactory interface {
    CreateSocket() (net.Conn, error)
}

pool, err := NewSocketPool(size, factory)
if err != nil {
    // Handle error
}
defer pool.Close()

conn, err := pool.Get()
if err != nil {
    // Handle error
}
defer conn.Close()

// Use conn to send and receive data
```

Here's an example Factory implementation that creates TCP connections:

```
type SocketFactory struct {
    addr string
}

func NewSocketFactory(addr string) SocketFactory {
    return SocketFactory{addr: addr}
}

func (factory SocketFactory) CreateSocket() (net.Conn, error) {
    return net.Dial("tcp", factory.addr)
}
```

You can also create a Factory that creates TLS connections:

```
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
```

## Test

First, You need to create a server private key and server certificate by using `gen-cert.sh`

Then, run the tests by using:

```
go test -v ./...
```

There is also a benchmark test that measures the performance:

```
go test ./... -bench=.
```

Benchmark result:

```
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkNormalTcpConnectionWithSocketPool-8              101617         11564 ns/op
BenchmarkNormalTcpConnectionWithoutSocketPool-8            58092         53496 ns/op
BenchmarkTlsConnectionWithSocketPool-8                     82755         12546 ns/op
BenchmarkTlsConnectionWithoutSocketPool-8                   1826        659079 ns/op
```
