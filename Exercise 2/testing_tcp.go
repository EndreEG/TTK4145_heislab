package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	SERVER = "your.server.ip.address"
	PORT   = "33546"
)

func main() {
	conn, err := net.Dial("tcp", SERVER+":"+PORT)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Read the welcome message
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}
	fmt.Println("Server welcome message:", string(buffer[:n]))

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		message = message + "\000" // Append null character

		// Send the message to the server
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		// Read the echo message
		n, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Println("Server echo message:", string(buffer[:n]))
	}
}
