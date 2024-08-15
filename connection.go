package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/bp-chat/bp-tui/commands"
)

type connection struct {
	conn        net.Conn
	reader      bufio.Reader
	writer      bufio.Writer
	receivedEof bool
	queue       CommandQueue
}

func (cnn *connection) IsOpen() bool {
	return cnn.receivedEof == false
}

func connect(address string) (*connection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &connection{
		conn,
		*bufio.NewReader(conn),
		*bufio.NewWriter(conn),
		false,
		CreateCommandQueue(),
	}, nil
}

func (cnn *connection) Send(outCommand commands.IOut) error {
	cmdId, err := cnn.queue.TakeSlot()
	if err != nil {
		cnn.queue.Enqueue(outCommand)
		return nil
	}
	cmd := outCommand.ToCommand(uint8(cmdId))
	data, err := cmd.Encode()
	if err != nil {
		return err
	}
	if _, err = cnn.writer.Write(data); err != nil {
		fmt.Printf("failed to send data: %x", data)
		return err
	}
	return cnn.writer.Flush()
}

func (cnn *connection) FreeCommandSlot(slot int) {
	cnn.queue.Free(slot)
}

func (cnn *connection) SendNext() error {
	out, err := cnn.queue.Pop()
	if err != nil {
		//the list is empty, there is no need to do anything special
		return nil
	}
	return cnn.Send(*out)
}

func (cnn *connection) Receive() (*commands.Command, error) {
	buffer := make([]byte, 4096)
	_, err := cnn.reader.Read(buffer)
	if err != nil {
		log.Printf("err: %v", err)
		cnn.receivedEof = true
		return nil, err
	}
	cmd, err := commands.Decode(buffer)
	if err != nil {
		log.Printf("could not parse data\n%s\n", err)
		return nil, err
	}
	return cmd, nil
}

func (cnn *connection) Close() {
	cnn.conn.Close()
}
