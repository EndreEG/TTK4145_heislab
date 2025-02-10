package main

//FOR WINDOWS

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	heartbeatInterval = 1 * time.Second // Time between heartbeats
	missedThreshold   = 3               // Number of missed heartbeats before taking over
	tcpPort           = ":9999"         // TCP port for communication
	counterFile       = "counter.txt"   // File to store the last number
)

// Read the last number from the file
func readLastNumber() int {
	data, err := os.ReadFile(counterFile)
	if err != nil {
		return 0 // Start from 0 if file doesn't exist
	}
	number, _ := strconv.Atoi(string(data))
	return number
}

// Write the last number to the file
func writeLastNumber(number int) {
	os.WriteFile(counterFile, []byte(strconv.Itoa(number)), 0644)
}

// Primary process: counts and sends heartbeats to the backup
func primary() {
	number := readLastNumber()
	fmt.Println("Primary is running. Starting count from:", number)

	// Set up TCP listener
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		fmt.Println("Error setting up TCP listener:", err)
		return
	}
	defer listener.Close()

	// Accept a connection from the backup
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	for {
		// Print the number and save it
		fmt.Println(number)
		writeLastNumber(number)
		number++

		// Send a heartbeat to the backup
		_, err := conn.Write([]byte("heartbeat\n"))
		if err != nil {
			fmt.Println("Error sending heartbeat:", err)
			return
		}

		// Wait before the next iteration
		time.Sleep(heartbeatInterval)
	}
}

// Backup process: connects to the primary and monitors for heartbeats
func backup() {
	fmt.Println("Backup is running. Waiting to take over...")

	// Connect to the primary
	conn, err := net.Dial("tcp", "localhost"+tcpPort)
	if err != nil {
		fmt.Println("Error connecting to primary:", err)
		fmt.Println("Assuming primary is dead. Taking over...")
		primary() // Become the primary
		return
	}
	defer conn.Close()

	missedHeartbeats := 0
	buffer := make([]byte, 1024)

	for {
		// Set a timeout for receiving heartbeats
		conn.SetReadDeadline(time.Now().Add(heartbeatInterval * 2))

		// Try to read a heartbeat
		_, err := conn.Read(buffer)
		if err != nil {
			missedHeartbeats++
			fmt.Println("Missed heartbeat:", missedHeartbeats)
			if missedHeartbeats >= missedThreshold {
				fmt.Println("Primary is dead. Taking over...")
				primary() // Become the primary
				return
			}
		} else {
			missedHeartbeats = 0 // Reset counter if heartbeat is received
		}
	}
}

func main() {
	// Check if this process is the primary or backup
	if len(os.Args) > 1 && os.Args[1] == "--primary" {
		primary()
	} else {
		backup()
	}
}
