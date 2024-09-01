package client

import "github.com/bp-chat/bp-tui/commands"

type EphemeralUser struct {
	Name      commands.UserName
	Keys      KeySet
	SharedKey [32]byte
	IsKeySet  bool
}
