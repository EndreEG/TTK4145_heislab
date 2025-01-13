package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Select an option:")
	fmt.Println("1. Connect to the server (TCP Client)")
	fmt.Println("2. Send 'Connect to' message")
	fmt.Println("3. Listen for incoming connections (TCP Server)")
	fmt.Print("Enter your choice: ")

	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1:
		connectToServer()
	case 2:
		sendConnectToMessage()
	case 3:
		startServer()
	default:
		fmt.Println("Invalid choice")
	}
}

func connectToServer() {
	serverIP := "server-ip-here" // Replace with actual server IP
	fixedSizePort := 34933
	delimitedPort := 33546

	fmt.Println("Choose the port to connect:")
	fmt.Println("1. Fixed-size messages (34933)")
	fmt.Println("2. Delimited messages (33546)")
	fmt.Print("Enter your choice: ")

	var portChoice int
	fmt.Scan(&portChoice)

	port := fixedSizePort
	if portChoice == 2 {
		port = delimitedPort
	}

	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to the server.")

	// Handle incoming messages
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\0') // Reading till '\0'
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}
			fmt.Println("Server:", strings.TrimRight(message, "\x00"))
		}
	}()

	// Send messages to the server
	sendMessages(conn)
}

func sendConnectToMessage() {
	serverIP := "server-ip-here" // Replace with actual server IP
	serverPort := 33546
	listenPort := 40000 // Port for your program to listen for incoming connections

	// Sending the "Connect to" message to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localIP := getLocalIP()
	connectMessage := fmt.Sprintf("Connect to: %s:%d\x00", localIP, listenPort)
	_, err = conn.Write([]byte(connectMessage))
	if err != nil {
		fmt.Println("Error sending connect message:", err)
	}
	fmt.Println("Sent connect message to server.")

	// Start listening for incoming connections
	startServer()
}

func startServer() {
	listenPort := 40000
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Listening for incoming connections on port %d...\n", listenPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func sendMessages(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if len(message) == 0 {
			continue
		}

		// Append null terminator for delimited messages
		if _, err := conn.Write([]byte(message + "\x00")); err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	panic("No local IP found")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		message, err := reader.ReadString('\0')
		if err != nil {
			fmt.Println("Connection closed by client.")
			return
		}

		fmt.Printf("Received: %s\n", strings.TrimRight(message, "\x00"))
		_, _ = writer.WriteString("Echo: " + message)
		writer.Flush()
	}
}
