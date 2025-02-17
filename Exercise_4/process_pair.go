package main2

//kill -INT $(lsof -t -i :8080)

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
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
	spawnBackup()

	count := startNum
	var conn net.Conn

	for {
		if conn == nil {
			fmt.Println("Waiting for backup...")
			conn, err = ln.Accept()
			if err != nil {
				fmt.Println("Failed to accept connection:", err)
				continue
			}
			fmt.Println(count)
		}

		writer := bufio.NewWriter(conn)
		_, err := fmt.Fprintln(writer, count)
		if err != nil {
			fmt.Println("Lost connection to backup. Waiting for new backup...")
			conn.Close()
			conn = nil
			spawnBackup()
			continue
		}

		writer.Flush()
		fmt.Println(count)
		count++
		time.Sleep(heartbeatFreq)
	}
}

func runBackup() {

	signal.Ignore(syscall.SIGTERM)
	
	for {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Println("Primary not found. Becoming new primary...")
			runPrimary(1) 
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

func spawnBackup() {
	cmd := exec.Command(os.Args[0], "backup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Start()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "backup" {
		runBackup()
	} else {
		runPrimary(1)
	}
}
