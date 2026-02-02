package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)


func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  Subscribe: go run . sub [topic]")
		fmt.Println("  Publish:   go run . pub [topic] [message]")
		return
	}

	mode := os.Args[1]  
	topic := os.Args[2] 

	conn, err := net.Dial("tcp", "localhost:9092")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	if mode == "sub" {
		cmd := Packet{OpCode: SUB, Topic: topic}
		WritePacket(conn, cmd)

		fmt.Printf("Listening for messages on topic: [%s]...\n", topic)
		for {
			p, err := ReadPacket(conn)
			if err != nil { return }
			fmt.Printf("received > %s\n", string(p.Payload))
		}

	} else if mode == "pub" {
		msg := strings.Join(os.Args[3:], " ")
		
		cmd := Packet{OpCode: PUB, Topic: topic, Payload: []byte(msg)}
		WritePacket(conn, cmd)
		fmt.Println("Message sent!")
	}
}