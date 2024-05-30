package main

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/bp-chat/bp-tui/commands"
	"github.com/bp-chat/bp-tui/commands/calls"
)

type CommandCode uint16

const (
	CNN CommandCode = 0
	MSG             = 1
)

const Host string = "127.0.0.1:6680"

func main() {
	log.Printf("trying to connect to %s...\n", Host)
	conn, err := connect(Host)
	reader := bufio.NewReader(os.Stdin)
	if err != nil {
		log.Fatalf("Could not connect to %s\n%s\n", Host, err)
	}
	defer conn.Close()
	log.Printf("Connected to %s...\n", Host)
	go listen(conn)
	for {
		cmdMsg := getMessage(reader)
		if err = conn.Send(cmdMsg.ToCommand()); err != nil {
			log.Fatalf("Could not write to connection\n%s\n", err)
		}
	}
}

func listen(cnn *connection) {
	for cnn.IsOpen() {
		cnn.Receive()
		time.Sleep(2 * time.Second)
	}
}

func getMessage(reader *bufio.Reader) calls.Message {
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Could not parse input\n%s\n", err)
	}
	return calls.Message{
		Header: commands.Header{
			Version: 0,
			SyncId:  0,
			Id:      0,
		},
		Recipient: "self",
		Message:   input,
	}
}
