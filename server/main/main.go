package main

import (
	"context"
	"net"
	"os"
	hw "server/hello_world"
	"time"

	"fmt"

	capnp "capnproto.org/go/capnp/v3"
	rpc "capnproto.org/go/capnp/v3/rpc"
)

// Timeout for server to respond
const deadline = 5 * time.Second

// Main function
func main() {
	server := hw.GreeterServer{}
	// capnp client capability that can be shared over the network
	client := hw.Greeter_ServerToClient(server)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	d := time.Now().Add(deadline)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	defer cancel()
	defer ln.Close()
	fmt.Println("Listening on 8080")
	for {
		// Listen for an incoming connection.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(ctx, client, conn)
	}
}

func handleRequest(ctx context.Context, client hw.Greeter, rwc net.Conn) error {
	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{
		BootstrapClient: capnp.Client(client),
	})
	defer conn.Close()
	select {
	case <-conn.Done():
		return nil
	case <-ctx.Done():
		return conn.Close()
	}
}
