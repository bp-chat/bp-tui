package commands

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type IOut interface {
	ToCommand(syncId uint8) Command
}

type Command struct {
	Header
	Args [][]uint8
}

const separator uint8 = 0x0a

func (cmd *Command) Encode() ([]uint8, error) {
	buffer := new(bytes.Buffer)

	if err := write(buffer, cmd.Version, false); err != nil {
		return nil, err
	}

	if err := write(buffer, cmd.SyncId, false); err != nil {
		return nil, err
	}
	if err := write(buffer, cmd.Id, false); err != nil {
		return nil, err
	}
	for i, element := range cmd.Args {
		if err := write(buffer, element, i+1 == len(cmd.Args)); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func write(b *bytes.Buffer, data any, is_last bool) error {
	var err = binary.Write(b, binary.BigEndian, data)
	if err != nil {
		return err
	}
	if is_last {
		return binary.Write(b, binary.BigEndian, separator)
	}
	return binary.Write(b, binary.BigEndian, separator)
}

func Decode(source []uint8, arg_count int) (*Command, error) {
	if len(source) < 8 {
		return nil, errors.New("source is to small, it must be at least 8 bytes")
	}
	if arg_count > 10 {
		return nil, errors.New("It looks like to many args")
	}
	if source[2] != separator || source[4] != separator || source[7] != separator {
		return nil, errors.New("Invalid data format")
	}
	var args [][]uint8
	var c_arg []uint8
	for i := 8; i < len(source); i++ {
		if len(args) == arg_count {
			break
		}

		elem := source[i]
		if elem != separator {
			c_arg = append(c_arg, elem)
			if i < len(source)-1 {
				continue
			}
		}
		args = append(args, c_arg)
		c_arg = make([]uint8, 0)
	}

	if len(args) != arg_count {
		msg := fmt.Sprintf(
			"Argument count mismatch\nexpected:%v and found %v",
			arg_count,
			len(args))
		return nil, errors.New(msg)
	}
	commandId := uint16(source[5])<<8 | uint16(source[6])
	header := Header{
		uint16(source[0])<<8 | uint16(source[1]),
		source[3],
		CommandCode(commandId),
	}

	return &Command{
		header,
		args,
	}, nil
}
