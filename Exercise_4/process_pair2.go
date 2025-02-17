package main

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
	address       = "localhost:8082"
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

	// Handle Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Primary received interrupt, exiting...")
		os.Exit(0) // Only exit primary, backup remains
	}()

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
	signal.Ignore(syscall.SIGTERM) // Ignore SIGTERM to avoid shutdown when primary exits

	for {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Println("Primary not found. Becoming new primary...")
			runPrimary(0) // Start from last known value (or 0 if no state was saved)
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

			lastCount, _ = strconv.Atoi(message[:len(message)-1])
			conn.SetReadDeadline(time.Now().Add(timeout))
		}
	}
}

func spawnBackup() {
	cmd := exec.Command(os.Args[0], "backup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // Detach backup from primary
	cmd.Start()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "backup" {
		runBackup()
	} else {
		runPrimary(0)
	}
}
