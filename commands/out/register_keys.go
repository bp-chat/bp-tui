package out

import "github.com/bp-chat/bp-tui/commands"

type RegisterKeys struct {
	IdKey     []byte
	SignedKey []byte
	Signature []byte
}

// ToCommand implements commands.IOut.
func (r RegisterKeys) ToCommand(header commands.Header) commands.Command {
	return commands.Command{
		Header: header,
		Args: [][]byte{
			r.IdKey,
			r.SignedKey,
			r.SignedKey,
		},
	}
}
