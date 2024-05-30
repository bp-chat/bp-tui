package calls

import "github.com/bp-chat/bp-tui/commands"

type Message struct {
	Header    commands.Header
	Recipient string
	Message   string
}

func (m *Message) ToCommand() commands.Command {
	header := commands.Header{
		Version: m.Header.Version,
		SyncId:  m.Header.SyncId,
		Id:      m.Header.Id,
	}
	args := [][]uint8{
		[]uint8(m.Recipient),
		[]uint8(m.Message),
	}

	return commands.Command{
		Header: header,
		Args:   args,
	}
}

func FromCommand(cmd *commands.Command) Message {
	headers := cmd.Header
	r := "whoops"
	m := "failed to read message"
	if cmd.Args[0] != nil {
		r = string(cmd.Args[0])
	}
	if cmd.Args[1] != nil {
		m = string(cmd.Args[1])
	}
	return Message{
		Header:    headers,
		Recipient: r,
		Message:   m,
	}
}
