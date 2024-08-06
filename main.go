package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bp-chat/bp-tui/commands"
	"github.com/bp-chat/bp-tui/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type ephemeralUser struct {
	name      commands.UserName
	keys      KeySet
	sharedKey [32]byte
	isKeySet  bool
}

const Host string = "127.0.0.1:6680"

func main() {
	fmt.Printf("\n\nWho are you\n")
	reader := bufio.NewReader(os.Stdin)
	name := []byte(getMessage(reader))
	var username commands.UserName
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

	p := tea.NewProgram(ui.New(func(nm string) error {
		return send(conn, &eu, nm)
	}, func() {
		broadcastCommand(conn)
	}))
	go listen(conn, p, &eu)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func listen(cnn *connection, teaProgam *tea.Program, eu *ephemeralUser) {
	for cnn.IsOpen() {
		bpMsg, err := cnn.Receive()
		if err != nil {
			teaProgam.Send(err)
			break
		}
		switch bpMsg.Header.Id {
		case commands.RKS:
			other, err := commands.NewRegisterKeys(bpMsg.Body)
			if err != nil {
				teaProgam.Send(err)
			}
			otherKeys := PublicKeySet{
				identityKey:  other.IdKey[:],
				signedKey:    other.SignedKey[:],
				ephemeralkey: other.EphemeralKey[:],
				signature:    other.Signature[:],
			}
			eu.sharedKey, err = CreateSharedKey(eu.keys, otherKeys)
			if err != nil {
				teaProgam.Send(err)
				break
			}
			eu.isKeySet = true
			break
		case commands.MSG:
			if eu.isKeySet == false {
				teaProgam.Send(errors.New("shared was not created"))
				break
			}
			encryptedMessage, err := commands.NewMessage(bpMsg.Body)
			if err != nil {
				teaProgam.Send(err)
				break
			}
			decrypted, err := Decrypt(eu.sharedKey[:], encryptedMessage.InitialVector, encryptedMessage.Message[:encryptedMessage.Len])
			if err != nil {
				teaProgam.Send(err)
				break
			}
			m := ui.Message{
				From:    string(encryptedMessage.Recipient[:]),
				Message: string(decrypted),
			}
			teaProgam.Send(m)
			break
		}

		// time.Sleep(1500 * time.Millisecond)
	}
}

func send(cnn *connection, user *ephemeralUser, textMsg string) error {
	if user.isKeySet == false {
		return errors.New("Shared key not setted yeat")
	}
	msgBytes := []byte(textMsg)
	if len(msgBytes) > commands.MessageSize {
		return errors.New("The message is to large")
	}
	iv, encrypted, err := Encrypt(user.sharedKey[:], msgBytes)
	if err != nil {
		return err
	}
	mlen := len(encrypted)
	if mlen > commands.MessageSize {
		return errors.New("The encrypted message is to large")
	}
	msgb := make([]byte, commands.MessageSize)
	copy(msgb, encrypted)
	msg := commands.Message{
		Recipient:     user.name,
		InitialVector: iv,
		Len:           int32(mlen),
		Message:       [commands.MessageSize]byte(msgb),
	}
	return cnn.Send(msg)
}

func broadcastCommand(cnn *connection) error {
	cmd := commands.BroadcastKeys{}
	return cnn.Send(cmd)
}

func registerE2eeKeys(cnn *connection, user ephemeralUser) error {
	cmd := commands.RegisterKeys{
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
