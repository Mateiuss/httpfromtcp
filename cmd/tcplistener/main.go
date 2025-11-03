package main

import (
	"fmt"
	"io"
	"strings"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer f.Close()
		defer close(out)

		curr := ""

		for {
			b := make([]byte, 8)

			_, err := f.Read(b)

			if err == io.EOF {
				break
			}

			s := string(b)
			splits := strings.Split(s, "\n")

			if len(splits) == 1 {
				curr += splits[0]
			} else {
				out <- curr + splits[0]
				curr = splits[1]
			}
		}

		if len(curr) != 0 {
			out <- curr
		}
	}()

	return out
}

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

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}

		fmt.Println("connection closed")
	}
}
