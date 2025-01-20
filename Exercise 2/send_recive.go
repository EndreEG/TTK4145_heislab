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

	//Get workspace number and calculate communication port
	fmt.Print("Enter your workspace number: ")
	var workspaceNumber int
	fmt.Scan(&workspaceNumber)
	communicationPort := 20000 + workspaceNumber

	fmt.Printf("Server IP found: %s\n", serverIP)
	fmt.Printf("Using communication port: %d\n", communicationPort)

	//WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Start listening for replies from the server
	go func() {
		defer wg.Done()
		listenForReplies(communicationPort)
	}()

	//Start sending messages to the server
	go func() {
		defer wg.Done()
		sendMessages(serverIP, communicationPort)
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

func listenForReplies(port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("Listening for replies on port %d...\n", port)

	buffer := make([]byte, 1024)
	for {
		n, senderAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error receiving reply:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Reply from %s: %s\n", senderAddr, message)
	}
}

func sendMessages(serverIP string, port int) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", serverIP, port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message to send: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}
