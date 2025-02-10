package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "time"
)

const (
    heartbeatInterval = 2 * time.Second
    heartbeatTimeout  = 5 * time.Second
)

func main() {
    role := os.Getenv("ROLE")
    countStr := os.Getenv("COUNT")
    count := 1
    if countStr != "" {
        count, _ = strconv.Atoi(countStr)
    }

    if role == "backup" {
        runBackup(count)
    } else {
        runPrimary(count)
    }
}

func runPrimary(count int) {
    fmt.Println("Starting as primary with count:", count)

    // Start backup process
    backupCmd := exec.Command(os.Args[0])
    backupCmd.Env = append(os.Environ(), "ROLE=backup", fmt.Sprintf("COUNT=%d", count))
    backupIn, _ := backupCmd.StdinPipe()
    backupOut, _ := backupCmd.StdoutPipe()
    backupCmd.Start()

    backupReader := bufio.NewReader(backupOut)

    ticker := time.NewTicker(heartbeatInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Send heartbeat
            fmt.Fprintln(backupIn, "heartbeat", count)
            count++
            fmt.Println("Count:", count)
        default:
            // Check for backup acknowledgment
            backupCmd.Process.Signal(os.Interrupt)
            ack, err := backupReader.ReadString('\n')
            if err != nil {
                fmt.Println("Backup process failed. Exiting...")
                return
            }
            if ack == "ack\n" {
                // Backup acknowledged heartbeat
                continue
            }
        }
    }
}

func runBackup(count int) {
    fmt.Println("Starting as backup with count:", count)
    reader := bufio.NewReader(os.Stdin)
    lastHeartbeat := time.Now()

    for {
        // Check for heartbeat
        os.Stdout.Write([]byte("ack\n"))
        line, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading heartbeat:", err)
            return
        }
        if line[:9] == "heartbeat" {
            lastHeartbeat = time.Now()
            count, _ = strconv.Atoi(line[10:])
        }

        // Check for heartbeat timeout
        if time.Since(lastHeartbeat) > heartbeatTimeout {
            fmt.Println("Primary missed heartbeat. Promoting to primary...")
            runPrimary(count)
            return
        }
    }
}
