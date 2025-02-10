package main

import (
    "fmt"
    "net"
    "os"
    "os/exec"
    "strconv"
    "time"
)

const (
    heartbeatInterval = 1 * time.Second
    missedThreshold   = 3
)

var counter = 1

func main() {
    // Start the backup process
    cmd := exec.Command("go", "run", "backup.go", strconv.Itoa(counter))
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Start()
    if err != nil {
        fmt.Println("Error starting backup process:", err)
        return
    }

    // Listen on port 12345 for heartbeat messages
    listener, err := net.Listen("tcp", ":12345")
    if err != nil {
        fmt.Println("Error listening on port 12345:", err)
        return
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            return
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    for {
        // Send heartbeat message
        _, err := conn.Write([]byte("heartbeat"))
        if err != nil {
            fmt.Println("Error sending heartbeat:", err)
            return
        }

        // Print the counter value
        fmt.Println(counter)
        counter++

        time.Sleep(heartbeatInterval)
    }
}
