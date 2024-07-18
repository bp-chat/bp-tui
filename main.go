package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bp-chat/bp-tui/commands"
	"github.com/bp-chat/bp-tui/commands/out"
	"github.com/bp-chat/bp-tui/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type ephemeralUser struct {
	name commands.UserName
	keys KeySet
}

const Host string = "127.0.0.1:6680"

func main() {
	fmt.Printf("\n\nWho are you\n")
	reader := bufio.NewReader(os.Stdin)
	name := []byte(getMessage(reader))
	var username [16]byte
	copy(username[:], name[:])
	eu := ephemeralUser{
		name: username,
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
	msgBytes := []byte(textMsg)
	if len(msgBytes) > out.MessageSize {
		return errors.New("The message is to large")
	}
	msg := out.Message{
		Recipient: user.name,
		Message:   [out.MessageSize]byte(msgBytes),
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
		IdKey:        [32]byte(user.keys.publicKey),
		SignedKey:    [32]byte(user.keys.preKey.PublicKey().Bytes()),
		Signature:    [64]byte(user.keys.signature),
		EphemeralKey: [32]byte(user.keys.ek.PublicKey().Bytes()),
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
