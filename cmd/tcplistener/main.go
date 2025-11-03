package main

import (
	"fmt"
	"net"
	"thechallange/internal/request"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		return
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		fmt.Printf("accepted a connection from %s\n", conn.RemoteAddr().String())

		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("error while getting the request from the connection")
			continue
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		fmt.Println("connection closed")
	}
}
