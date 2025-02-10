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

func runPrimary(startNum int) {
	ln, err := net.Listen("tcp", address)
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
		}
	}
}

func spawnBackup(startNum int) {
	cmd := exec.Command(os.Args[0], "backup")
	cmd.Env = append(os.Environ(), fmt.Sprintf("START_NUM=%d", startNum))
	cmd.Start()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "backup" {
		runBackup()
	} else {
		runPrimary(1)
	}
}
