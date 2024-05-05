package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

type CommandCode uint16

const (
	CNN CommandCode = 0
	MSG             = 1
)

type Command struct {
	version uint16
	id      uint8
	code    CommandCode
}

const Host string = "127.0.0.1:5501"

func main() {
	log.Printf("trying to connect to %s...\n", Host)
	conn, err := net.Dial("tcp", Host)
	if err != nil {
		log.Fatalf("Could not connect to %s\n%s\n", Host, err)
	}
	cmd := Command{
		version: 0,
		id:      1,
		code:    CommandCode(CNN),
	}
	log.Printf("Connected to %s...\n", Host)
	buffer := new(bytes.Buffer)
	if err = binary.Write(buffer, binary.BigEndian, cmd); err != nil {
		log.Printf("Could not parse cmd %v\n%s\n", cmd, err)
	}
	if _, err = conn.Write(buffer.Bytes()); err != nil {
		log.Fatalf("Could not write to connection\n%s\n", err)
	}
	conn.Close()
}
