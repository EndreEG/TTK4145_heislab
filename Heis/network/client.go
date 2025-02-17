package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type ElevatorState struct {
	ID       string `json:"id"`
	Floor    int    `json:"floor"`
	Behavior string `json:"behavior"`
}

func StartClient(serverAddr, id string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Could not connect to primary:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to primary at", serverAddr)

	// Send elevator ID to server
	fmt.Fprintln(conn, id)

	writer := bufio.NewWriter(conn)

	for {
		state := ElevatorState{
			ID:       id,
			Floor:    getCurrentFloor(),
			Behavior: getCurrentBehavior(),
		}

		data, _ := json.Marshal(state)
		writer.WriteString(string(data) + "\n")
		writer.Flush()

		time.Sleep(1 * time.Second)
	}
}

func getCurrentFloor() int {
	// Replace with actual floor retrieval
	return 2
}

func getCurrentBehavior() string {
	// Replace with actual behavior retrieval
	return "MovingUp"
}
