package main

import (
	"errors"

	cl "github.com/bp-chat/bp-tui/client"
	"github.com/bp-chat/bp-tui/commands"
	"github.com/bp-chat/bp-tui/ui"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

const Host string = "127.0.0.1:6680"

func main() {

	name := []byte("test user")
	var username commands.UserName
	copy(username[:], name[:])
	eu := cl.EphemeralUser{
		Name: username,
		Keys: cl.CreateKeys(),
	}
	// log.Printf("trying to connect to %s...\n", Host)
	conn, err := cl.Connect(Host)
	if err != nil {
		log.Fatalf("Could not connect to %s\n%s\n", Host, err)
	}
	defer conn.Close()
	// log.Printf("Connected to %s...\n", Host)
	client := cl.New(eu, conn)
	client.RefreshKeys()
	p := tea.NewProgram(ui.New(client), tea.WithAltScreen())
	go listen(conn, p, &eu)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func listen(cnn *cl.Connection, teaProgam *tea.Program, eu *cl.EphemeralUser) {
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
			otherKeys := cl.PublicKeySet{
				IdentityKey:  other.IdKey[:],
				SignedKey:    other.SignedKey[:],
				Ephemeralkey: other.EphemeralKey[:],
				Signature:    other.Signature[:],
			}
			eu.SharedKey, err = cl.CreateSharedKey(eu.Keys, otherKeys)
			if err != nil {
				teaProgam.Send(err)
				break
			}
			eu.IsKeySet = true
			break
		case commands.MSG:
			if eu.IsKeySet == false {
				teaProgam.Send(errors.New("shared was not created"))
				break
			}
			encryptedMessage, err := commands.NewMessage(bpMsg.Body)
			if err != nil {
				teaProgam.Send(err)
				break
			}
			decrypted, err := cl.Decrypt(eu.SharedKey[:], encryptedMessage.InitialVector, encryptedMessage.Message[:encryptedMessage.Len])
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

func send(cnn *cl.Connection, user *cl.EphemeralUser, textMsg string) error {
	if user.IsKeySet == false {
		return errors.New("Shared key not setted yeat")
	}
	msgBytes := []byte(textMsg)
	if len(msgBytes) > commands.MessageSize {
		return errors.New("The message is to large")
	}
	iv, encrypted, err := cl.Encrypt(user.SharedKey[:], msgBytes)
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
		Recipient:     user.Name,
		InitialVector: iv,
		Len:           int32(mlen),
		Message:       [commands.MessageSize]byte(msgb),
	}
	return cnn.Send(msg)
}

func broadcastCommand(cnn *cl.Connection) error {
	cmd := commands.BroadcastKeys{}
	return cnn.Send(cmd)
}

func registerE2eeKeys(cnn *cl.Connection, user cl.EphemeralUser) error {
	cmd := commands.RegisterKeys{
		User:         user.Name,
		IdKey:        [32]byte(user.Keys.PublicKey),
		SignedKey:    [32]byte(user.Keys.PreKey.PublicKey().Bytes()),
		Signature:    [64]byte(user.Keys.Signature),
		EphemeralKey: [32]byte(user.Keys.Ek.PublicKey().Bytes()),
	}
	return cnn.Send(cmd)
}
