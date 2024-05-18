package commands

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Command struct {
	version uint16
	sync_id uint8
	id      uint16
	args    [][]uint8
}

func (cmd *Command) encode() ([]uint8, error) {
	buffer := new(bytes.Buffer)
	err_msg := errors.New("Command could not be parsed")

	if err := write(buffer, cmd.version); err != nil {
		return nil, err_msg
	}
	if err := write(buffer, 0x0a); err != nil {
		return nil, err_msg
	}
	return buffer.Bytes(), nil
}

func write(b *bytes.Buffer, data any) error {
	return binary.Write(b, binary.BigEndian, data)
}

func decode(data []uint8) Command {
	return Command{
		0,
		1,
		2,
		[][]uint8{},
	}
}
