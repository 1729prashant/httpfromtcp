package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	filePath := "messages.txt"

	fileToRead, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fileToRead.Close()

	scanner := bufio.NewScanner(fileToRead)

	// Iterate over each line in the file.
	for scanner.Scan() {
		line := scanner.Text()         // Get the current line as a string.
		fmt.Printf("read: %s\n", line) // Process or print the line.
	}

	// Check for any errors encountered during scanning.
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
