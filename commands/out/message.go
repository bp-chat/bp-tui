package out

import (
	"errors"

	"github.com/bp-chat/bp-tui/commands"
)

type Message struct {
	header    commands.Header
	Recipient string
	Message   string
}

func (m Message) ToCommand(syncId uint8) commands.Command {
	args := [][]uint8{
		[]uint8(m.Recipient),
		[]uint8(m.Message),
	}

	return commands.Command{
		Header: commands.NewHeader(commands.MSG, syncId),
		Args:   args,
	}
}

func FromCommand(cmd *commands.Command) (*Message, error) {
	headers := cmd.Header
	if len(cmd.Args) != 2 {
		return nil, errors.New("invalid argument count")
	}
	if cmd.Args[0] == nil || cmd.Args[1] == nil {
		return nil, errors.New("could not parse arguments")
	}
	return &Message{
		headers,
		string(cmd.Args[0]),
		string(cmd.Args[1]),
	}, nil
}
