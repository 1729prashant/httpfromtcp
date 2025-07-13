package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Resolve the address we want to send UDP packets to: localhost:42069
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("error resolving address: %v", err)
	}

	// Establish a UDP connection to that address
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("error dialing UDP: %v", err)
	}
	defer conn.Close() // Ensure the connection is closed on exit

	// Create a reader that reads from standard input (the terminal)
	reader := bufio.NewReader(os.Stdin)

	for {
		// Prompt the user for input
		fmt.Print("> ")

		// Read a line from stdin (up to and including the newline character)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error reading from stdin: %v", err)
			continue
		}

		// Send the line over the UDP connection
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("error sending UDP packet: %v", err)
		}
	}
}
