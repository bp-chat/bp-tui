package main

import (
	"errors"

	"github.com/bp-chat/bp-tui/commands"
)

type Client struct {
	user ephemeralUser
	conn *connection
}

func (client *Client) SendMessage(msg string) error {
	user := client.user
	if user.isKeySet == false {
		return errors.New("Shared key not setted yeat")
	}
	msgBytes := []byte(msg)
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
	msgBuffer := make([]byte, commands.MessageSize)
	copy(msgBuffer, encrypted)
	return client.conn.Send(
		commands.Message{
			Recipient:     user.name,
			InitialVector: iv,
			Len:           int32(mlen),
			Message:       [commands.MessageSize]byte(msgBuffer),
		})
}

func (client *Client) RefreshKeys() error {
	user := client.user
	cmd := commands.RegisterKeys{
		User:         user.name,
		IdKey:        [32]byte(user.keys.publicKey),
		SignedKey:    [32]byte(user.keys.preKey.PublicKey().Bytes()),
		Signature:    [64]byte(user.keys.signature),
		EphemeralKey: [32]byte(user.keys.ek.PublicKey().Bytes()),
	}
	return client.conn.Send(cmd)
}

func (client *Client) BroadcastKeys() error {
	cmd := commands.BroadcastKeys{}
	return client.conn.Send(cmd)
}
