package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
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
}
