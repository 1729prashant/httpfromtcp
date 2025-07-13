package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

		lines := getLinesChannel(conn)

		// Loop over the channel and print each line as it comes in
		for line := range lines {
			fmt.Println(line)
		}

		// When the channel is closed (EOF or client disconnect), we log that
		fmt.Println("connection closed...")

	}

}

// getLinesChannel reads lines from an io.ReadCloser (like net.Conn),
// sends each line on a channel, and closes both the channel and connection when done.
func getLinesChannel(f io.ReadCloser) <-chan string {
	// Create an unbuffered channel of strings to hold the lines
	lines := make(chan string)

	// Start a new goroutine so that reading doesn't block the caller
	go func() {
		// Ensure that the file is closed once the goroutine is done reading
		defer f.Close()

		// Ensure that the channel is closed after all lines are sent
		defer close(lines)

		// Create a scanner that reads from the file line by line.
		scanner := bufio.NewScanner(f)

		// Loop over each line in the file.
		// scanner.Scan() returns true until the end of the file or an error occurs.
		for scanner.Scan() {
			// scanner.Text() returns the current line as a string (without newline).
			// Send this line to the channel so the main function can consume it.
			lines <- scanner.Text()
		}

		// After the loop, check whether an error occurred during scanning.
		// EOF is not treated as an error, but I/O problems (e.g., disk issues) are.
		if err := scanner.Err(); err != nil {
			// Print the error to standard error (not standard output)
			fmt.Fprintln(os.Stderr, "Error reading from connection:", err)
		}
	}()

	// Return the channel immediately so the caller can start consuming lines
	return lines
}
