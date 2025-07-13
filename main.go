package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	filePath := "messages.txt"
	buffer8bytes := make([]byte, 8)

	fileToRead, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fileToRead.Close()

	for {
		n, err := fileToRead.Read(buffer8bytes)

		if err != nil {
			if err == io.EOF {
				fmt.Println("End of file reached.")
			} else {
				fmt.Println("Error reading file:", err)
			}
			break // Exit the loop on error or EOF
		}

		// Process the read bytes (n bytes were read into buffer[:n])
		//fmt.Printf("Read %d bytes: %s (raw bytes: %v)\n", n, string(buffer8bytes[:n]), buffer8bytes[:n])
		fmt.Printf("read: %s\n", string(buffer8bytes[:n]))
	}
}
