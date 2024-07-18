package commands

type BroadcastKeys struct {
}

// ToCommand implements commands.IOut.
func (m BroadcastKeys) ToCommand(syncId uint8) Command {
	args := []uint8{}

	return Command{
		Header: NewHeader(BKS, syncId),
		Body:   args,
	}
}
