package out

import "github.com/bp-chat/bp-tui/commands"

type BroadcastKeys struct {
}

// ToCommand implements commands.IOut.
func (m BroadcastKeys) ToCommand(syncId uint8) commands.Command {
	args := [][]uint8{}

	return commands.Command{
		Header: commands.NewHeader(commands.BKS, syncId),
		Args:   args,
	}
}
