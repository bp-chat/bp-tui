package commands

import (
	"bytes"
	"encoding/binary"
)

const MessageSize = 980

type Message struct {
	Recipient     UserName
	InitialVector [12]byte //TODO READ IV SIZE from e2ee
	Len           int32
	Message       [MessageSize]byte
}

func (m Message) ToCommand(syncId uint8) Command {
	buffer := new(bytes.Buffer)
	_ = binary.Write(buffer, binary.BigEndian, m)

	return Command{
		Header: NewHeader(MSG, syncId),
		Body:   buffer.Bytes(),
	}
}

func NewMessage(body []byte) (*Message, error) {
	var msg Message
	reader := bytes.NewReader(body)
	if err := binary.Read(reader, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
