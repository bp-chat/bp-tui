package out

import "github.com/bp-chat/bp-tui/commands"

type RegisterKeys struct {
	IdKey     []byte
	SignedKey []byte
	Signature []byte
}

// ToCommand implements commands.IOut.
func (r RegisterKeys) ToCommand(syncId uint8) commands.Command {
	return commands.Command{
		Header: commands.NewHeader(commands.RKS, syncId),
		Args: [][]byte{
			r.IdKey,
			r.SignedKey,
			r.SignedKey,
		},
	}
}
