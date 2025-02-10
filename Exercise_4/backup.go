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

var counter int

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run backup.go <starting_counter>")
        return
    }
    counter, _ = strconv.Atoi(os.Args[1])

    // Listen on port 12346 for heartbeat messages from the primary
    listener, err := net.Listen("tcp", ":12346")
    if err != nil {
        fmt.Println("Error listening on port 12346:", err)
        return
    }
    defer listener.Close()

    missedCount := 0

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            return
        }

        go handleConnection(conn, &missedCount)
    }
}

func handleConnection(conn net.Conn, missedCount *int) {
    defer conn.Close()

    buf := make([]byte, 1024)

    for {
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Error reading from connection:", err)
            return
        }

        if n > 0 && string(buf[:n]) == "heartbeat" {
            *missedCount = 0
        } else {
            *missedCount++
        }

        if *missedCount >= missedThreshold {
            // Become the new primary
            fmt.Println("Becoming the new primary")

            // Start the new backup process
            cmd := exec.Command("go", "run", "backup.go", strconv.Itoa(counter))
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            err = cmd.Start()
            if err != nil {
                fmt.Println("Error starting new backup process:", err)
                return
            }

            // Start counting from where the previous primary left off
            for {
                fmt.Println(counter)
                counter++

                time.Sleep(heartbeatInterval)
            }
        }
    }
}
