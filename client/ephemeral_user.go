package client

import "github.com/bp-chat/bp-tui/commands"

type EphemeralUser struct {
	name      commands.UserName
	keys      KeySet
	sharedKey [32]byte
	isKeySet  bool
}
