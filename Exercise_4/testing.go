package main

//FOR WINDOWS

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

const (
	HOST           = "localhost"
	PORT           = "5000"
	HEARTBEAT_RATE = 500 * time.Millisecond
	FAILOVER_TIME  = 2 * HEARTBEAT_RATE
)

var count int = 1
var mutex sync.Mutex

func StartBackup() {
	cmd := exec.Command(os.Args[0], "backup", strconv.Itoa(count))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start backup:", err)
		os.Exit(1)
	}
}


func RunPrimary() {
	StartBackup()

	listener, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error starting listener:", err)
		os.Exit(1)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)

	for {
		mutex.Lock()
		fmt.Println(count)
		count++
		mutex.Unlock()

		_, err := writer.WriteString(fmt.Sprintf("%d\n", count))
		if err != nil {
			fmt.Println("Lost connection to backup. Exiting...")
			break
		}
		writer.Flush()

		time.Sleep(HEARTBEAT_RATE)
	}
}


func RunBackup(initialCount int) {
	count = initialCount
	conn, err := net.Dial("tcp", HOST+":"+PORT)
	if err == nil {
		// Primary exists, act as backup
		defer conn.Close()
		reader := bufio.NewScanner(conn)
		lastHeartbeat := time.Now()

		for reader.Scan() {
			mutex.Lock()
			count, _ = strconv.Atoi(reader.Text())
			mutex.Unlock()
			lastHeartbeat = time.Now()
		}

		for time.Since(lastHeartbeat) < FAILOVER_TIME {
			time.Sleep(HEARTBEAT_RATE)
		}

		fmt.Println("Primary seems dead. Taking over...")
	} else {
		fmt.Println("No primary detected. Becoming primary.")
	}

	RunPrimary()
}

func main() {
	if len(os.Args) > 2 && os.Args[1] == "backup" {
		initialCount, _ := strconv.Atoi(os.Args[2])
		RunBackup(initialCount)
	} else {
		RunPrimary()
	}
}
