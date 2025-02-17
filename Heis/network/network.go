package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
<<<<<<< HEAD
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
=======
)

type ElevatorState struct {
	ID       string `json:"id"`
	Floor    int    `json:"floor"`
	Behavior string `json:"behavior"`
}

// Send elevator state update to the primary server
func SendElevatorUpdate(id string, floor int, behavior string) {
	message := ElevatorState{ID: id, Floor: floor, Behavior: behavior}
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	conn, err := net.Dial("tcp", "localhost:5000") // Connect to the server
	if err != nil {
		fmt.Println("Could not connect to server:", err)
		return
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	writer.WriteString(string(data) + "\n")
	writer.Flush()
>>>>>>> 971c211... Implemented network functionality
}
