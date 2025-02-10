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

	go startListeningServer(":40000")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		listenForRepliesTCP(conn)
	}()

	go func() {
		defer wg.Done()
		if strings.ToLower(mode) == "1" {
			sendFixedSizeMessages(conn)
		} else {
			sendNullTerminatedMessages(conn)
		}
	}()

	wg.Wait()
}

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

		_, err = conn.Write([]byte("Echo: " + message))
		if err != nil {
			fmt.Println("Error sending echo:", err)
			return
		}
	}
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
