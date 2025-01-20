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
	// Step 1: Start a listener for incoming connections
	go startListeningServer(":40000") // Listens on port 40000 for incoming connections

	// Step 2: Find the server's IP by listening on port 30000
	serverIP := findServerIP()

	// Step 3: Ask for the message mode (fixed-size or null-terminated)
	fmt.Print("Enter message mode (fixed-size or null-terminated): ")
	var mode string
	fmt.Scan(&mode)

	var tcpPort int
	switch strings.ToLower(mode) {
	case "fixed-size":
		tcpPort = 34933
	case "null-terminated":
		tcpPort = 33546
	default:
		fmt.Println("Invalid mode. Use 'fixed-size' or 'null-terminated'.")
		return
	}

	fmt.Printf("Server IP found: %s\n", serverIP)
	fmt.Printf("Connecting to TCP port: %d\n", tcpPort)

	// Step 4: Establish a TCP connection
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, tcpPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server.")

	// Step 5: Create WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Step 6: Start listening for replies from the server
	go func() {
		defer wg.Done()
		listenForRepliesTCP(conn)
	}()

	// Step 7: Start sending messages to the server
	go func() {
		defer wg.Done()
		if strings.ToLower(mode) == "fixed-size" {
			sendFixedSizeMessages(conn)
		} else {
			sendNullTerminatedMessages(conn)
		}
	}()

	// Wait for goroutines to finish (they won't unless there's an error)
	wg.Wait()
}

// Step 1: Start Listening Server
func startListeningServer(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting listening server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening for incoming connections on %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle each connection in a separate goroutine
		go handleIncomingConnection(conn)
	}
}

func handleIncomingConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by peer:", err)
			return
		}
		fmt.Printf("Message received: %s\n", strings.TrimSpace(message))

		// Echo back the message (optional)
		_, err = conn.Write([]byte("Echo: " + message))
		if err != nil {
			fmt.Println("Error sending echo:", err)
			return
		}
	}
}

// Other functions (findServerIP, listenForRepliesTCP, sendFixedSizeMessages, sendNullTerminatedMessages)
// remain unchanged and are included below for reference.

func findServerIP() string {
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

		paddedMessage := fmt.Sprintf("%-256s", message)
		_, err := conn.Write([]byte(paddedMessage))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}

func sendNullTerminatedMessages(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter null-terminated message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		messageWithNull := message + "\x00"
		_, err := conn.Write([]byte(messageWithNull))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}
