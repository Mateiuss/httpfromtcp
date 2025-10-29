package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err1 := reader.ReadString('\n')
		if err1 != nil {
			fmt.Println("an error occured")
			continue
		}

		_, err2 := conn.Write([]byte(line))

		if err2 != nil {
			fmt.Println("an error occured while sending the line")
		}
	}
}
