package out

import (
	"fmt"

	"github.com/bp-chat/bp-tui/commands"
)

type RegisterKeys struct {
	User         string
	IdKey        []byte
	SignedKey    []byte
	Signature    []byte
	EphemeralKey []byte //will be removed later when we implement a direct way of communication
	EphemeralSig []byte //this too
}

// ToCommand implements commands.IOut.
func (r RegisterKeys) ToCommand(syncId uint8) commands.Command {
	fmt.Printf("\n keys to command parse")
	fmt.Printf("\n idkey: %x", r.IdKey)
	return commands.Command{
		Header: commands.NewHeader(commands.RKS, syncId),
		Args: [][]byte{
			[]byte(r.User),
			r.IdKey,
			r.SignedKey,
			r.Signature,
			r.EphemeralKey,
			r.EphemeralSig,
		},
	}
}
