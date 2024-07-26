package commands

import (
	"bytes"
	"encoding/binary"
)

const RegisterKeysSize = 176

type RegisterKeys struct {
	User         UserName
	IdKey        [32]byte
	SignedKey    [32]byte
	Signature    [64]byte
	EphemeralKey [32]byte //will be removed later when we implement a direct way of communication
}

// ToCommand implements commands.IOut.
func (r RegisterKeys) ToCommand(syncId uint8) Command {
	buffer := new(bytes.Buffer)
	_ = binary.Write(buffer, binary.BigEndian, r)

	return Command{
		Header: NewHeader(RKS, syncId),
		Body:   buffer.Bytes(),
	}
}

func NewRegisterKeys(body []byte) (*RegisterKeys, error) {
	var rks RegisterKeys
	reader := bytes.NewReader(body)
	if err := binary.Read(reader, binary.BigEndian, &rks); err != nil {
		return nil, err
	}
	return &rks, nil
}
