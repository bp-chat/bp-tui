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

const new_line uint8 = 0x0a

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

func (cmd *Command) EncodeSized() ([]uint8, error) {
	buffer := new(bytes.Buffer)

	if err := write_prepend(buffer, cmd.Version, false); err != nil {
		return nil, err
	}

	if err := write_prepend(buffer, cmd.SyncId, false); err != nil {
		return nil, err
	}
	if err := write_prepend(buffer, cmd.Id, false); err != nil {
		return nil, err
	}
	for _, element := range cmd.Args {
		if err := write_prepend(buffer, element, true); err != nil {
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
		return binary.Write(b, binary.BigEndian, new_line)
	}
	return binary.Write(b, binary.BigEndian, new_line)
}

func write_prepend(b *bytes.Buffer, data any, prepend bool) error {
	if prepend {
		s := uint32(binary.Size(data))
		err := binary.Write(b, binary.BigEndian, s)
		if err != nil {
			return err
		}
		return binary.Write(b, binary.BigEndian, data)
	}
	return binary.Write(b, binary.BigEndian, data)
}

func Decode(source []uint8, argCount int) (*Command, error) {
	if len(source) < 8 {
		return nil, errors.New("source is to small, it must be at least 8 bytes")
	}
	if argCount > 10 {
		return nil, errors.New("It looks like to many args")
	}
	if source[2] != new_line || source[4] != new_line || source[7] != new_line {
		return nil, errors.New("Invalid data format")
	}
	var args [][]uint8
	var c_arg []uint8
	for i := 8; i < len(source); i++ {
		if len(args) == argCount {
			break
		}

		elem := source[i]
		if elem != new_line {
			c_arg = append(c_arg, elem)
			if i < len(source)-1 {
				continue
			}
		}
		args = append(args, c_arg)
		c_arg = make([]uint8, 0)
	}

	if len(args) != argCount {
		msg := fmt.Sprintf(
			"Argument count mismatch\nexpected:%v and found %v",
			argCount,
			len(args))
		return nil, errors.New(msg)
	}
	commandId := binary.BigEndian.Uint16(source[5:7])
	header := Header{
		binary.BigEndian.Uint16(source[0:2]),
		source[3],
		CommandCode(commandId),
	}

	return &Command{
		header,
		args,
	}, nil
}

func DecodeSized(source []uint8, argCount int) (*Command, error) {
	if len(source) < 5 {
		return nil, errors.New("source is to small, it must be at least 5 bytes")
	}
	if argCount > 10 {
		return nil, errors.New("It looks like to many args")
	}
	var args [][]uint8
	var i uint32
	var lenSize uint32 = 4

	for i = 5; i < uint32(len(source)); {
		valueIdx := i + lenSize
		size := binary.BigEndian.Uint32(source[i:valueIdx])
		endIdx := valueIdx + size
		v := source[valueIdx:endIdx]
		args = append(args, v)
		i = endIdx
	}

	if len(args) != argCount {
		msg := fmt.Sprintf(
			"Argument count mismatch\nexpected:%v and found %v",
			argCount,
			len(args))
		return nil, errors.New(msg)
	}
	commandId := binary.BigEndian.Uint16(source[3:5])
	header := Header{
		binary.BigEndian.Uint16(source[0:2]),
		source[2],
		CommandCode(commandId),
	}

	return &Command{
		header,
		args,
	}, nil
}
