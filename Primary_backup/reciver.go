// coordinator/tcpserver.go
package main

import (
	"encoding/json"
	"fmt"
	"net" // Adjust if needed based on your project structure
	"bufio"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)






func main() {
	
	conn, err := net.Dial("tcp", primaryPort)
	if err != nil {
		fmt.Println("No primary found, starting as primary.")
		go startPrimary(primaryPort, secondaryPort)
	} else {
		conn.Close()
		fmt.Println("Primary found, starting as secondary.")
		go startSecondary(secondaryPort, primaryPort)
	}

	select {}
}





// Elevator state struct
const NumFloors int = 4
const NumButtons int = 3

type State struct {
	Elevator_id      int
	Elevator_floor   int
	Elevator_dir     int
	Elevator_behaviour int
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








//Primary and secondary server
const (
	primaryPort   = ":5000"
	secondaryPort = ":5001"
	stateFile     = "count_state.txt"
)

func loadCount() int {
	file, err := os.Open(stateFile)
	if err != nil {
		return 0 
	}
	defer file.Close()

	var count int
	fmt.Fscanf(file, "%d", &count)
	return count
}

func saveCount(count int) {
	file, err := os.Create(stateFile)
	if err != nil {
		fmt.Println("Error saving state:", err)
		return
	}
	defer file.Close()
	fmt.Fprintf(file, "%d", count)
}

func startPrimary(primaryAddr, secondaryAddr string) {
	listener, err := net.Listen("tcp", primaryAddr)
	if err != nil {
		fmt.Println("Error starting primary:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Primary started on", primaryAddr)

	go startSecondary(secondaryAddr, primaryAddr)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	count := loadCount()
	var conn net.Conn
	for {
		if conn == nil {
			fmt.Println("Connecting to secondary...")
			conn, err = net.Dial("tcp", secondaryAddr)
			if err != nil {
				fmt.Println("No secondary found. Retrying...")
				time.Sleep(2 * time.Second)
				continue
			}
		}

		select {
		case <-signalChan:
			fmt.Println("Primary shutting down, saving state and notifying secondary...")
			saveCount(count)
			conn.Close()
			return
		default:
			time.Sleep(1 * time.Second)
			count++
			saveCount(count)
			fmt.Println("Count:", count)
			fmt.Fprintf(conn, "%d\n", count)
		}
	}
}

func startSecondary(secondaryAddr, primaryAddr string) {
	listener, err := net.Listen("tcp", secondaryAddr)
	if err != nil {
		fmt.Println("Error starting secondary:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Secondary started on", secondaryAddr)

	count := loadCount()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		fmt.Println("Connected to primary")

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			count, _ = strconv.Atoi(scanner.Text())
			saveCount(count)
		}

		fmt.Println("Primary disconnected. Becoming primary...")
		conn.Close()
		go startPrimary(secondaryAddr, primaryAddr)
		return
	}
}
