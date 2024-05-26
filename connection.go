package main

import (
	"github.com/bp-chat/bp-tui/commands"
	"net"
)

type connection struct {
	conn net.Conn
}

func connect(address string) (*connection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &connection{
		conn,
	}, nil
}

func (cnn *connection) Send(cmd commands.Command) error {
	return nil
}

func (cnn *connection) Close() {
	cnn.conn.Close()
}
