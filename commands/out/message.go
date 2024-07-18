package out

import (
	"bytes"
	"encoding/binary"

	"github.com/bp-chat/bp-tui/commands"
)

const MessageSize = 1024

type Message struct {
	Recipient commands.UserName
	Message   [MessageSize]byte
}

func (m Message) ToCommand(syncId uint8) commands.Command {
	buffer := new(bytes.Buffer)
	_ = binary.Write(buffer, binary.BigEndian, m)

	return commands.Command{
		Header: commands.NewHeader(commands.MSG, syncId),
		Body:   buffer.Bytes(),
	}
}

func FromCommand(cmd *commands.Command) (*Message, error) {
	var msg Message
	reader := bytes.NewReader(cmd.Body)
	if err := binary.Read(reader, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
