package main

import (
	"encoding/binary"
	"io"
	"net"
)

const (
	PUB = 1
	SUB = 2
)

type Packet struct {
	OpCode  byte
	Topic   string
	Payload []byte
}
//sender to receiver (takes a packet and turns into raw data to send over network) - writing data
func WritePacket(conn net.Conn, cmd Packet) error {
	conn.Write([]byte{cmd.OpCode}) // 1 byte

	topicLen := uint16(len(cmd.Topic))//2 bytes (and using int 16)
	binary.Write(conn, binary.LittleEndian, topicLen)

	conn.Write([]byte(cmd.Topic))

	if cmd.OpCode == PUB {
		payloadLen := uint32(len(cmd.Payload))
		binary.Write(conn, binary.LittleEndian, payloadLen)
		conn.Write(cmd.Payload)
	}
	return nil
}
// Reading Data 
func ReadPacket(conn net.Conn) (Packet, error) {
	var p Packet

	op := make([]byte, 1)
	_, err := io.ReadFull(conn, op)
	if err != nil {
		return p, err
	}
	p.OpCode = op[0]

	var topicLen uint16
	if err := binary.Read(conn, binary.LittleEndian, &topicLen); err != nil {
		return p, err
	}

	topicBytes := make([]byte, topicLen)
	if _, err := io.ReadFull(conn, topicBytes); err != nil {
		return p, err
	}
	p.Topic = string(topicBytes)

	if p.OpCode == PUB {
		var payloadLen uint32
		if err := binary.Read(conn, binary.LittleEndian, &payloadLen); err != nil {
			return p, err
		}
		payloadBytes := make([]byte, payloadLen)
		if _, err := io.ReadFull(conn, payloadBytes); err != nil {
			return p, err
		}
		p.Payload = payloadBytes
	}

	return p, nil
}
