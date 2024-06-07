package main

import (
	"bufio"
	"github.com/bp-chat/bp-tui/commands"
	"log"
	"net"
)

type connection struct {
	conn        net.Conn
	reader      bufio.Reader
	writer      bufio.Writer
	receivedEof bool
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
	}, nil
}

func (cnn *connection) Send(cmd commands.Command) error {
	data, err := cmd.Encode()
	if err != nil {
		return err
	}
	if _, err = cnn.writer.Write(data); err != nil {
		return err
	}
	return cnn.writer.Flush()
}

func (cnn *connection) Receive() (*commands.Command, error) {
	buffer := make([]byte, 64)
	_, err := cnn.reader.Read(buffer)
	if err != nil {
		log.Printf("err: %v", err)
		cnn.receivedEof = true
		return nil, err
	}
	cmd, err := commands.Decode(buffer, 2)
	if err != nil {
		log.Printf("could not parse data\n%s\n", err)
		return nil, err
	}
	return cmd, nil
}

func (cnn *connection) Close() {
	cnn.conn.Close()
}
