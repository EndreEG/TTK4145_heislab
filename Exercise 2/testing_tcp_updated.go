package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

func main() {
	serverIP := findServerIP()

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

	// Start listening for incoming connections
	localPort := 40000
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		fmt.Println("Error setting up listener:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening for incoming connections on port %d\n", localPort)

	// Send the "Connect to" message to the server
	localIP, err := getLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		return
	}

	connectMessage := fmt.Sprintf("Connect to: %s:%d\x00", localIP, localPort)
	_, err = conn.Write([]byte(connectMessage))
	if err != nil {
		fmt.Println("Error sending 'Connect to' message:", err)
		return
	}

	fmt.Println("Sent 'Connect to' message to server.")

	// WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Handle replies from the original connection
	go func() {
		defer wg.Done()
		listenForRepliesTCP(conn)
	}()

	// Handle new incoming connections from the server
	go func() {
		defer wg.Done()
		for {
			newConn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting new connection:", err)
				continue
			}
			fmt.Println("Accepted new connection from server.")

			// Start a goroutine to handle this new connection
			go handleNewConnection(newConn)
		}
	}()

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
		return senderAddr.IP.String()
	}
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func listenForRepliesTCP(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading reply from server:", err)
			return
		}
		reply := string(buffer[:n])
		fmt.Printf("Reply from server: %s\n", strings.TrimSpace(reply))
	}
}

func handleNewConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from new connection:", err)
			return
		}

		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("Received from new connection: %s\n", message)

		// Echo the message back
		_, err = conn.Write([]byte("Echo: " + message))
		if err != nil {
			fmt.Println("Error writing to new connection:", err)
			return
		}
	}
}
