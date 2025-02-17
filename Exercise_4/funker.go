package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	primaryPort   = ":5000"
	secondaryPort = ":5001"
	stateFile     = "count_state.txt"
)

func main() {
	// Try to connect to an existing primary
	conn, err := net.Dial("tcp", primaryPort)
	if err != nil {
		// No primary found, become the primary
		fmt.Println("No primary found, starting as primary.")
		go startPrimary(primaryPort, secondaryPort)
	} else {
		// Primary exists, become the secondary
		conn.Close()
		fmt.Println("Primary found, starting as secondary.")
		go startSecondary(secondaryPort, primaryPort)
	}

	select {} // Keep main running
}

func loadCount() int {
	file, err := os.Open(stateFile)
	if err != nil {
		return 0 // Default to 0 if file doesn't exist
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

	// Ensure a secondary is started
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
			fmt.Println("Primary count:", count)
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
			fmt.Println("Received count from primary:", count)
		}

		fmt.Println("Primary disconnected. Becoming primary...")
		conn.Close()
		go startPrimary(secondaryAddr, primaryAddr) // Swap roles
		return
	}
}
