// Data Packet: "SEQ|DATA" → "3|hello"
// ACK Packet: "ACK:SEQ" → "ACK:3"

package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)
func main(){
	
	remoteAddr, err := net.ResolveUDPAddr("udp", "localhost:12345")
	if err != nil {
		fmt.Println("Error resolving address: ", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		fmt.Println("Error dialing: ", err)
		os.Exit(1)
	}

	messages := []string{
		"0|hello",
		"1|from",
		"2|zaher",
	}
	
	ackChan := make(chan string)
	go listenAck(conn, ackChan)
	
	timeout := 2 * time.Second

	for i, msg := range messages {
		for {
			fmt.Printf("Sending packet #%d...\n", i)
			conn.Write([]byte(msg))

			select{
				case ack := <- ackChan:
					if ack == fmt.Sprintf("ACK:%d", i) {
						fmt.Printf("ACK received for #%d\n", i)
						break
					}
				case <-time.After(timeout):
					fmt.Printf("Timeout on packet #%d, retrying...\n", i)
                	continue
			}
			break
		}
	}
}

func listenAck(conn *net.UDPConn, ch chan string){
	buf := make([]byte, 1024)
	for{
		n, _, err := conn.ReadFromUDP(buf)
		if err == nil {
			ack := string(buf[:n])
			if strings.HasPrefix(ack, "ACK:"){
				ch <- ack
			}
		}
	}
}