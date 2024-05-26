package main

import (
	"log"

	"github.com/bp-chat/bp-tui/commands"
	"github.com/bp-chat/bp-tui/commands/calls"
)

type CommandCode uint16

const (
	CNN CommandCode = 0
	MSG             = 1
)

const Host string = "127.0.0.1:5501"

func main() {
	log.Printf("trying to connect to %s...\n", Host)
	conn, err := connect(Host)
	if err != nil {
		log.Fatalf("Could not connect to %s\n%s\n", Host, err)
	}
	defer conn.Close()
	log.Printf("Connected to %s...\n", Host)
	cmdMsg := calls.Message{
		Header: commands.Header{
			Version: 0,
			SyncId:  0,
			Id:      0,
		},
		Recipient: "self",
		Message:   "hello world",
	}

	if err = conn.Send(cmdMsg.ToCommand()); err != nil {
		log.Fatalf("Could not write to connection\n%s\n", err)
	}
}
