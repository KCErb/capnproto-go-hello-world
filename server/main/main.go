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

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Println("Listening on 8080")
	for {
		// Listen for an incoming connection.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		d := time.Now().Add(deadline)
		ctx, _ := context.WithDeadline(context.Background(), d)
		// Handle connections in a new goroutine.
		go handleRequest(ctx, server, conn)
	}
}

type ErrorReporter struct{}

func (ErrorReporter) ReportError(err error) {
	fmt.Println(err)
}

func handleRequest(ctx context.Context, server hw.GreeterServer, rwc net.Conn) error {
	fmt.Println("creating client")
	client := hw.Greeter_ServerToClient(server)
	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{
		BootstrapClient: capnp.Client(client),
		ErrorReporter:   ErrorReporter{},
	})
	defer conn.Close()
	select {
	case <-conn.Done():
		fmt.Println("conn.Done")
		return nil
	case <-ctx.Done():
		fmt.Println("ctx.Cancel")
		return conn.Close()
	}
}
