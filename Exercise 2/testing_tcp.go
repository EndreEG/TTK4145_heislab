package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	//Find the server's IP
	serverIP := findServerIP()

	// Ask for the message mode (fixed size or null terminated)
	fmt.Print("Enter message mode (fixed size (1) or null terminated (2)): ")
	var mode string
	fmt.Scan(&mode)

	var tcpPort int
	switch strings.ToLower(mode) {
	case "1":
		tcpPort = 34933
	case "2":
		tcpPort = 33546
	default:
		fmt.Println("Invalid mode. Use '1' or '2'.")
		return
	}

	fmt.Printf("Server IP found: %s\n", serverIP)
	fmt.Printf("Connecting to TCP port: %d\n", tcpPort)

	// Establish a TCP connection
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, tcpPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server.")

	//WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	//Listening for replies from the server
	go func() {
		defer wg.Done()
		listenForRepliesTCP(conn)
	}()

	//Sending messages to the server
	go func() {
		defer wg.Done()
		if strings.ToLower(mode) == "1" {
			sendFixedSizeMessages(conn)
		} else {
			sendNullTerminatedMessages(conn)
		}
	}()

	// Wait for goroutines to finish (they won't unless there's an error)
	wg.Wait()
}

func findServerIP() string {
	// Create a UDP socket to listen on port 30000
	addr, err := net.ResolveUDPAddr("udp", ":30000")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Listening for server broadcasts on port 30000...")

	buffer := make([]byte, 1024)
	for {
		n, senderAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error receiving broadcast:", err)
			continue
		}

		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("Received broadcast from %s: %s\n", senderAddr, message)
		return senderAddr.IP.String() // Return the server's IP
	}
}

func listenForRepliesTCP(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		reply, err := reader.ReadBytes('\n') // Assumes server replies are newline-terminated
		if err != nil {
			fmt.Println("Error reading reply from server:", err)
			return
		}
		fmt.Printf("Reply from server: %s\n", strings.TrimSpace(string(reply)))
	}
}

func sendFixedSizeMessages(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter fixed-size message (up to 256 bytes): ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if len(message) > 256 {
			fmt.Println("Message too long. Please limit to 256 bytes.")
			continue
		}

		// Pad message to 256 bytes
		paddedMessage := fmt.Sprintf("%-256s", message)
		_, err := conn.Write([]byte(paddedMessage))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}

func sendNullTerminatedMessages(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter null-terminated message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		// Append null terminator
		messageWithNull := message + "\x00"
		_, err := conn.Write([]byte(messageWithNull))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}
