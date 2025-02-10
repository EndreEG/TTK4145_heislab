package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	address       = "localhost:8080"
	heartbeatFreq = 500 * time.Millisecond
	timeout       = 2 * heartbeatFreq
)

<<<<<<< HEAD

func readLastNumber() int {
	data, err := os.ReadFile(counterFile)
	if err != nil {
		return 0 
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
=======
func runPrimary(startNum int) {
	ln, err := net.Listen("tcp", address)
>>>>>>> 277fb674ec679e28a813aff71d18f6332c5fa9ba
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	// Spawn backup process
	spawnBackup(startNum)

	count := startNum
	for {
		fmt.Println(count)

		// Accept backup connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		// Send heartbeat with the current count
		writer := bufio.NewWriter(conn)
		fmt.Fprintln(writer, count)
		writer.Flush()
		conn.Close()

		count++
		time.Sleep(heartbeatFreq)
	}
}

func runBackup() {
	for {
		conn, err := net.Dial("tcp", address)
		if err != nil {
<<<<<<< HEAD
			missedHeartbeats++
			fmt.Println("Missed heartbeat:", missedHeartbeats)
			if missedHeartbeats >= missedThreshold {
				fmt.Println("Primary is dead. Taking over...")
				primary() 
				return
			}
		} else {
			missedHeartbeats = 0 
=======
			fmt.Println("Primary not found. Becoming new primary...")
			runPrimary(1) // Default to 1 if no previous count is known
			return
		}

		reader := bufio.NewReader(conn)
		conn.SetReadDeadline(time.Now().Add(timeout))

		var lastCount int
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Primary lost. Taking over...")
				runPrimary(lastCount + 1)
				return
			}

			// Update last received count
			lastCount, _ = strconv.Atoi(message[:len(message)-1])
			conn.SetReadDeadline(time.Now().Add(timeout))
>>>>>>> 277fb674ec679e28a813aff71d18f6332c5fa9ba
		}
	}
}

func spawnBackup(startNum int) {
	cmd := exec.Command(os.Args[0], "backup")
	cmd.Env = append(os.Environ(), fmt.Sprintf("START_NUM=%d", startNum))
	cmd.Start()
}

func main() {
<<<<<<< HEAD
	if len(os.Args) > 1 && os.Args[1] == "--primary" {
		primary()
=======
	if len(os.Args) > 1 && os.Args[1] == "backup" {
		runBackup()
>>>>>>> 277fb674ec679e28a813aff71d18f6332c5fa9ba
	} else {
		runPrimary(1)
	}
}
