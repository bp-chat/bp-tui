package main

import (
	"bufio"
	"log"
	"time"

	"github.com/bp-chat/bp-tui/commands"
	"github.com/bp-chat/bp-tui/commands/calls"
	"github.com/bp-chat/bp-tui/ui"
	tea "github.com/charmbracelet/bubbletea"
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
	if err != nil {
		log.Fatalf("Could not connect to %s\n%s\n", Host, err)
	}
	defer conn.Close()
	log.Printf("Connected to %s...\n", Host)

	p := tea.NewProgram(ui.New(func(nm string) {
		send(conn, nm)
	}))
	go listen(conn, p)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func listen(cnn *connection, teaProgam *tea.Program) {
	for cnn.IsOpen() {
		bpMsg, err := cnn.Receive()
		if err != nil {
			teaProgam.Send(err)
		} else {
			teaProgam.Send(bpMsg)
		}
		time.Sleep(2 * time.Second)
	}
}

func send(cnn *connection, msg string) {
	cmsg := calls.Message{
		Header: commands.Header{
			Version: 0,
			SyncId:  0,
			Id:      0,
		},
		Recipient: "self",
		Message:   msg,
	}
	cnn.Send(cmsg.ToCommand())
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
