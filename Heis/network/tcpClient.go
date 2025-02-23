package network

import (
	"encoding/json"
	"log"
	"net"
)

func SendElevatorState(state interface{}) {
	conn, err := net.Dial("tcp", "localhost:3000") // Change to your primary’s IP and port
	if err != nil {
		log.Println("Failed to connect to primary:", err)
		return
	}
	defer conn.Close()

	data, err := json.Marshal(state)
	if err != nil {
		log.Println("Failed to encode state:", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Println("Failed to send state update:", err)
	}
}
