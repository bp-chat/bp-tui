package main

import (
	"errors"

	cl "github.com/bp-chat/bp-tui/client"
	"github.com/bp-chat/bp-tui/commands"
	"github.com/bp-chat/bp-tui/ui"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

const host string = "127.0.0.1:6680"

func main() {
	//TODO add logs
	config := cl.Config{
		Host: host,
	}
	p := tea.NewProgram(ui.New(config), tea.WithAltScreen())
	// go listen(conn, p, &eu)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// TODO move most of this logic to the client?
// maybe init the client inside the tea init method?
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
	}
}
