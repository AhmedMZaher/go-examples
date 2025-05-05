package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	remoteAddr, err := net.ResolveUDPAddr("udp", "localhost:12345")
	if err != nil {
		fmt.Println("Error resolving address: ", err)
		os.Exit(1)
	}

    conn, _ := net.ListenUDP("udp", remoteAddr)
    defer conn.Close()

    buffer := make([]byte, 1024)
    expectedSeq := 0
    packetMap := make(map[int]string)

    for {
        n, addr, _ := conn.ReadFromUDP(buffer)
        msg := string(buffer[:n])

        parts := strings.SplitN(msg, "|", 2)
        seq, _ := strconv.Atoi(parts[0])
        data := parts[1]

        fmt.Printf("Received #%d: %s\n", seq, data)
        packetMap[seq] = data

        ack := fmt.Sprintf("ACK:%d", seq)
        conn.WriteToUDP([]byte(ack), addr)
        fmt.Printf("Sent %s\n", ack)

        for {
            data, ok := packetMap[expectedSeq]
            if !ok {
                break
            }
            fmt.Printf("Processing #%d: %s\n", expectedSeq, data)
            delete(packetMap, expectedSeq)
            expectedSeq++
        }

        fmt.Println("Buffered:", sortedKeys(packetMap))
    }
}

func sortedKeys(m map[int]string) []int {
    keys := []int{}
    for k := range m {
        keys = append(keys, k)
    }
    sort.Ints(keys)
    return keys
}
