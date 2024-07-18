package out

import (
	"bytes"
	"encoding/binary"

	"github.com/bp-chat/bp-tui/commands"
)

type RegisterKeys struct {
	User         commands.UserName
	IdKey        [32]byte
	SignedKey    [32]byte
	Signature    [64]byte
	EphemeralKey [32]byte //will be removed later when we implement a direct way of communication
}

// ToCommand implements commands.IOut.
func (r RegisterKeys) ToCommand(syncId uint8) commands.Command {
	buffer := new(bytes.Buffer)
	_ = binary.Write(buffer, binary.BigEndian, r)

	return commands.Command{
		Header: commands.NewHeader(commands.RKS, syncId),
		Body:   buffer.Bytes(),
	}
}
