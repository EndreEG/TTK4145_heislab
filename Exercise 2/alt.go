package main

import (
    "fmt"
    "net"
    "os"
    "time"
)

func main() {
    // Channel to communicate server IP address
    serverIPChan := make(chan string)

    // Goroutine to receive UDP packets and find the server IP
    go func() {
        addr := net.UDPAddr{
            Port: 30000,
            IP:   net.ParseIP("0.0.0.0"),
        }

        conn, err := net.ListenUDP("udp", &addr)
        if err != nil {
            fmt.Println("Error: ", err)
            os.Exit(1)
        }
        defer conn.Close()

        buffer := make([]byte, 1024)

        for {
            numBytesReceived, fromAddr, err := conn.ReadFromUDP(buffer)
            if err != nil {
                fmt.Println("Error: ", err)
                continue
            }

            fmt.Printf("Received %d bytes from %s: %s\n", numBytesReceived, fromAddr.String(), string(buffer[:numBytesReceived]))
            fmt.Println("Server IP address: ", fromAddr.IP.String())
            serverIPChan <- fromAddr.IP.String()
            break
        }
    }()

    // Wait for the server IP address
    serverIP := <-serverIPChan

    // Sending UDP packets and receiving responses
    sendUDPPackets(serverIP, 20000, "Hello, server!")
}

func sendUDPPackets(serverIP string, port int, message string) {
    addr := net.UDPAddr{
        Port: port,
        IP:   net.ParseIP(serverIP),
    }

    conn, err := net.DialUDP("udp", nil, &addr)
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(1)
    }
    defer conn.Close()

    _, err = conn.Write([]byte(message))
    if err != nil {
        fmt.Println("Error: ", err)
    }

    buffer := make([]byte, 1024)
    conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Set a read deadline
    numBytesReceived, fromAddr, err := conn.ReadFromUDP(buffer)
    if err != nil {
        fmt.Println("Error: ", err)
        return
    }

    fmt.Printf("Received %d bytes from %s: %s\n", numBytesReceived, fromAddr.String(), string(buffer[:numBytesReceived]))
}
