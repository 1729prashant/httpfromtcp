package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	filePath := "messages.txt"

	fileToRead, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Pass the opened file to getLinesChannel to get a channel of lines
	lines := getLinesChannel(fileToRead)

	// Loop over the lines received from the channel.
	// The loop automatically stops when the channel is closed.
	for line := range lines {
		// Print each line with a "read: " prefix.
		// This keeps output behavior consistent with the original code.
		fmt.Printf("read: %s\n", line)
	}
}

// getLinesChannel starts a goroutine that reads lines from the given file,
// sends each line over a channel, and closes both the file and channel when done.
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
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
		}
	}()

	// Return the channel immediately so the caller can start consuming lines
	return lines
}
