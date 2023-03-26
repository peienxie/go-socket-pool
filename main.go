package main

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/peienxie/socketpool/echo"
	"github.com/peienxie/socketpool/factory"
	"github.com/peienxie/socketpool/pool"
)

func main() {
	listener, err := echo.StartEchoServer("127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	pool, err := pool.NewSocketPool(8, factory.NewSocketFactory(listener.Addr().String()))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			conn, err := pool.Get()
			if err != nil {
				log.Printf("Failed to get socket connection from pool: %v", err)
			} else {
				defer pool.Release(conn)

				msg := fmt.Sprintf("hello from connection %d\n", id)
				conn.Write([]byte(msg))
				log.Print("Sent: " + msg)

				resp, err := io.ReadAll(conn)
				if err == nil {
					log.Print("Received: " + string(resp))
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
