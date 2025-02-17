package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

var (
	connections    = make(map[string]net.Conn)
	mu             sync.Mutex
	elevatorStates = make(map[string]ElevatorState) // Stores elevator states
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

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Lost connection")
			return
		}

		var state ElevatorState
		err = json.Unmarshal([]byte(message), &state)
		if err != nil {
			fmt.Println("Invalid message:", message)
			continue
		}

		mu.Lock()
		elevatorStates[state.ID] = state
		mu.Unlock()

		fmt.Printf("Received update: %+v\n", state)
	}
}
