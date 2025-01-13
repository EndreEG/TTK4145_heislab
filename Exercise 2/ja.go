package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// TCP server addresses
	fixedSizeAddr := "localhost:34933"  // Port for fixed-size messages
	delimitedAddr := "localhost:33546"  // Port for delimited messages

	// Connect to both servers
	fixedSizeConn, err := net.Dial("tcp", fixedSizeAddr)
	if err != nil {
		fmt.Println("Error connecting to fixed-size server:", err)
		return
	}
	defer fixedSizeConn.Close()
	fmt.Println("Connected to fixed-size server at", fixedSizeAddr)

	delimitedConn, err := net.Dial("tcp", delimitedAddr)
	if err != nil {
		fmt.Println("Error connecting to delimited server:", err)
		return
	}
	defer delimitedConn.Close()
	fmt.Println("Connected to delimited server at", delimitedAddr)

	// Set TCP_NODELAY to disable Nagle's algorithm (prevent coalescing)
	setTCPNoDelay(fixedSizeConn)
	setTCPNoDelay(delimitedConn)

	// Start receiving from both servers in parallel
	go receiveFixedSizeMessages(fixedSizeConn)
	go receiveDelimitedMessages(delimitedConn)

	// Start sending messages
	go sendMessages(fixedSizeConn, "fixed-size")
	go sendMessages(delimitedConn, "delimited")

	// Wait for user input to terminate
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Set the TCP_NODELAY socket option to prevent packet coalescing
func setTCPNoDelay(conn net.Conn) {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		err := tcpConn.SetNoDelay(true)
		if err != nil {
			fmt.Println("Error setting TCP_NODELAY:", err)
		}
	}
}

// Receive and process fixed-size messages (1024 bytes) from the server
func receiveFixedSizeMessages(conn net.Conn) {
	buffer := make([]byte, 1024) // Buffer to hold exactly 1024 bytes
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading fixed-size message:", err)
			return
		}
		fmt.Printf("Received fixed-size message (%d bytes): %s\n", n, string(buffer[:n]))
	}
}

// Receive and process delimited messages (terminated with \0) from the server
func receiveDelimitedMessages(conn net.Conn) {
	for {
		message, err := readDelimitedMessage(conn)
		if err != nil {
			fmt.Println("Error reading delimited message:", err)
			return
		}
		fmt.Printf("Received delimited message: %s\n", message)
	}
}

// Read a message terminated by a null byte (\0) from the connection
func readDelimitedMessage(conn net.Conn) (string, error) {
	var message []byte
	buffer := make([]byte, 1)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return "", err
		}
		if n > 0 {
			if buffer[0] == '\x00' { // Check for the null-terminator
				break
			}
			message = append(message, buffer[0]) // Append character to the message
		}
	}

	return string(message), nil
}

// Send messages to the server based on the message type (fixed-size or delimited)
func sendMessages(conn net.Conn, messageType string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message) // Remove any trailing newline characters

		if messageType == "fixed-size" {
			// Send a fixed-size message (1024 bytes)
			fixedSizeMessage := message
			if len(fixedSizeMessage) < 1024 {
				fixedSizeMessage = fixedSizeMessage + strings.Repeat(" ", 1024-len(fixedSizeMessage)) // Pad message
			} else if len(fixedSizeMessage) > 1024 {
				fixedSizeMessage = fixedSizeMessage[:1024] // Trim message to 1024 bytes if too long
			}

			_, err := conn.Write([]byte(fixedSizeMessage))
			if err != nil {
				fmt.Println("Error sending fixed-size message:", err)
				return
			}
			fmt.Println("Sent fixed-size message:", fixedSizeMessage)
		} else if messageType == "delimited" {
			// Send a delimited message (append '\0')
			_, err := conn.Write([]byte(message + "\x00"))
			if err != nil {
				fmt.Println("Error sending delimited message:", err)
				return
			}
			fmt.Println("Sent delimited message:", message)
		}

		// Wait a bit before sending the next message
		time.Sleep(2 * time.Second)
	}
}
