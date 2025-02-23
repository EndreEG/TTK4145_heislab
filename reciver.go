// coordinator/tcpserver.go
package main

import (
	"encoding/json"
	"fmt"
	"net" // Adjust if needed based on your project structure
)

// Eleator state struct
const NumFloors int = 4
const NumButtons int = 3

type State struct {
	Elevator_id      int
	Elevator_floor   int
	Elevator_dir     int
	Elevator_request [NumFloors][NumButtons]int
}

// StartTCPServer listens for incoming elevator state updates
func StartTCPServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

// handleConnection decodes incoming elevator state and processes it
func handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var state State

	if err := decoder.Decode(&state); err != nil {
		fmt.Println("Failed to decode state:", err)
		return
	}

	fmt.Printf("Received state update: %+v\n", state)
	// Here you can add the state to some central data structure
}

func main() {
	StartTCPServer("3000") // Pick whatever port you like

}
