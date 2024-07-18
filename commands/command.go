package commands

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type IOut interface {
	ToCommand(syncId uint8) Command
}

type Command struct {
	Header
	Body []byte
}

type In = interface{}

func (cmd *Command) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.BigEndian, cmd.Header); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, cmd.Body); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Decode(source []byte) (*Command, error) {
	if len(source) < headerSize {
		return nil, errors.New("source is to small, it must be at least 8 bytes")
	}
	var header Header
	headerReader := bytes.NewReader(source[0:headerSize])
	if err := binary.Read(headerReader, binary.BigEndian, &header); err != nil {
		return nil, err
	}
	return &Command{
		header,
		source[headerSize:],
	}, nil
}

func (cmd *Command) Parse() (In, error) {
	switch cmd.Header.Id {
	case MSG:
		return NewMessage(cmd.Body)
	case BKS:
		return NewRegisterKeys(cmd.Body)
	default:
		return nil, errors.New("default")
	}
}
