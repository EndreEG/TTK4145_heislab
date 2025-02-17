package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type ElevatorState struct {
	ID       string `json:"id"`
	Floor    int    `json:"floor"`
	Behavior string `json:"behavior"`
}

var (
	connections = make(map[string]net.Conn)
	mu          sync.Mutex
)

func StartServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Read elevator ID
	id, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading ID:", err)
		return
	}

	id = id[:len(id)-1] // Trim newline
	mu.Lock()
	connections[id] = conn
	mu.Unlock()
	fmt.Println("Elevator connected:", id)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Lost connection to", id)
			mu.Lock()
			delete(connections, id)
			mu.Unlock()
			return
		}

		var state ElevatorState
		err = json.Unmarshal([]byte(message), &state)
		if err != nil {
			fmt.Println("Invalid message from", id, ":", message)
			continue
		}

		fmt.Printf("Received from %s: %+v\n", id, state)
	}
}
