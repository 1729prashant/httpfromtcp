package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	protocol := "tcp"
	port := ":42069"
	l, err := net.Listen(protocol, port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("connection established...")

		//lines := getLinesChannel(conn)
		req, err := request.RequestFromReader(conn)

		if err != nil {
			log.Println("Failed to parse request:", err)
			continue
		}

		// Print the parsed RequestLine in the specified format.
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		fmt.Println("connection closed...")

	}

}
