package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bp-chat/bp-tui/commands/out"
	"github.com/bp-chat/bp-tui/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type ephemeralUser struct {
	name string
	keys KeySet
}

const Host string = "127.0.0.1:6680"

func main() {
	fmt.Printf("\n\nWho are you\n")
	reader := bufio.NewReader(os.Stdin)
	name := getMessage(reader)
	eu := ephemeralUser{
		name: name,
		keys: CreateKeys(),
	}
	log.Printf("trying to connect to %s...\n", Host)
	conn, err := connect(Host)
	if err != nil {
		log.Fatalf("Could not connect to %s\n%s\n", Host, err)
	}
	defer conn.Close()
	log.Printf("Connected to %s...\n", Host)
	registerE2eeKeys(conn, eu)

	p := tea.NewProgram(ui.New(func(nm string) {
		send(conn, eu, nm)
	}, func() {
		broadcastCommand(conn)
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
		time.Sleep(1500 * time.Millisecond)
	}
}

func send(cnn *connection, user ephemeralUser, textMsg string) error {
	msg := out.Message{
		Recipient: user.name,
		Message:   textMsg,
	}
	return cnn.Send(msg)
}

func broadcastCommand(cnn *connection) error {
	cmd := out.BroadcastKeys{}
	return cnn.Send(cmd)
}

func registerE2eeKeys(cnn *connection, user ephemeralUser) error {
	cmd := out.RegisterKeys{
		User:         user.name,
		IdKey:        user.keys.publicKey,
		SignedKey:    user.keys.preKey.PublicKey().Bytes(),
		Signature:    user.keys.signature,
		EphemeralKey: user.keys.ek.PublicKey().Bytes(),
	}
	return cnn.Send(cmd)
}

func getMessage(reader *bufio.Reader) string {
	input, isPrefix, err := reader.ReadLine()
	if err != nil {
		log.Fatalf("Could not parse input\n%s\n", err)
	}
	if isPrefix {
		log.Fatalf("Use a smaller name my friend")
	}
	return string(input)
}
