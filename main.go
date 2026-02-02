package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type Broker struct {
	// A map to track who is listening to what.
	// Key: Topic Name (e.g., "sports")
	// Value: List of connections (people listening)
	subscribers map[string][]net.Conn
	
	mutex       sync.RWMutex
	walFile     *os.File
}

func main() {
	fmt.Println("Starting")

	f, err := os.OpenFile("wal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { panic(err) }

	broker := &Broker{
		subscribers: make(map[string][]net.Conn),
		walFile:     f,
	}

	listener, err := net.Listen("tcp", ":9092")
	if err != nil { panic(err) }
	defer listener.Close()
	
	fmt.Println("Listening on port 9092")

	for {
		conn, err := listener.Accept()
		if err != nil { continue }
		go broker.handleConnection(conn)
	}
}

func (b *Broker) handleConnection(conn net.Conn) {	
	for {
		packet, err := ReadPacket(conn)
		if err != nil {
			conn.Close() 
			return 
		}

		if packet.OpCode == PUB {
			b.handlePublish(packet)
		} else if packet.OpCode == SUB {
			b.handleSubscribe(conn, packet.Topic)
		}
	}
}

func (b *Broker) handlePublish(p Packet) {
	b.walFile.WriteString(fmt.Sprintf("%s|%s\n", p.Topic, string(p.Payload)))
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	subscribers := b.subscribers[p.Topic]
	
	fmt.Printf(" Received msg on [%s]. Forwarding to %d subscribers\n", p.Topic, len(subscribers))
	for _, sub := range subscribers {
		WritePacket(sub, p) 
	}
}

func (b *Broker) handleSubscribe(conn net.Conn, topic string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.subscribers[topic] = append(b.subscribers[topic], conn)
	fmt.Printf("👤 New subscriber connected to topic: [%s]\n", topic)
}