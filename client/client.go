package client

import (
	"errors"

	"github.com/bp-chat/bp-tui/commands"
)

type Client struct {
	user EphemeralUser
	conn *Connection
}

const MaxNumberOfCommands int = 16

func New(user EphemeralUser, conn *Connection) Client {
	return Client{
		user,
		conn,
	}
}

func (client *Client) SendMessage(msg string) error {
	user := client.user
	if user.IsKeySet == false {
		return errors.New("Shared key not setted yeat")
	}
	msgBytes := []byte(msg)
	if len(msgBytes) > commands.MessageSize {
		return errors.New("The message is to large")
	}
	iv, encrypted, err := Encrypt(user.SharedKey[:], msgBytes)
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
			Recipient:     user.Name,
			InitialVector: iv,
			Len:           int32(mlen),
			Message:       [commands.MessageSize]byte(msgBuffer),
		})
}

func (client *Client) RefreshKeys() error {
	user := client.user
	cmd := commands.RegisterKeys{
		User:         user.Name,
		IdKey:        [32]byte(user.Keys.PublicKey),
		SignedKey:    [32]byte(user.Keys.PreKey.PublicKey().Bytes()),
		Signature:    [64]byte(user.Keys.Signature),
		EphemeralKey: [32]byte(user.Keys.Ek.PublicKey().Bytes()),
	}
	return client.conn.Send(cmd)
}

func (client *Client) BroadcastKeys() error {
	cmd := commands.BroadcastKeys{}
	return client.conn.Send(cmd)
}
