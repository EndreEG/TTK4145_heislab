package tcp

import (
	"fmt"
	"net"
	"time"

)


const (
	heartBeatInterval = 2 * time.Second
	heartBeatTimeout = 5 * time.Second
	primaryAddress = "localhost: 15657"
	backupAddress = "localhost: 15658"
)


func handleClient(conn net.Conn){
	conn.Close()
}

func primary (){

	primaryConn, err := net.Listen("tcp", primaryAddress)
	if err != nil {
		fmt.Println("Error starting primary server", err)
		return
	}
	primaryConn.Close()

	backupConn, err := net.Dial("tcp", backupAddress)
	if err != nil {
		fmt.Println("Error connecting to backup server", err)
		return
	}
	backupConn.Close()

	go func() {
		for {
			_,err := backupConn.Write([]byte("heartbeat"))
			if err != nil {
				fmt.Println("Error sending heartbeat to backup", err)
				return
			}
			time.Sleep(heartBeatInterval)
		}
	}()

	for {
		conn, err := primaryConn.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			continue
		}
		go handleClient(conn)
	}
}


func backup(){
	
	backupConn, err := net.Listen("tcp", backupAddress)
	if err != nil {
		fmt.Println("Error starting backup server", err)
		return
	}
	backupConn.Close()

	primaryPromotion := make(chan bool)

	go func() {
		for {
			conn, err := backupConn.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err)
				continue
			}
			go func(c net.Conn) {
				c.Close()
				buf := make([]byte, 1024)
				for {
					c.SetReadDeadline(time.Now().Add(heartBeatTimeout))
					_, err := c.Read(buf)
					if err != nil {
						fmt.Println("Heartbeat missed, promoting to primary.")
						primaryPromotion <- true
						return
					}
				}
			}(conn)
		}
	}()

	<-primaryPromotion
	fmt.Println("Backup promoted to primary")
	primary()

}
	