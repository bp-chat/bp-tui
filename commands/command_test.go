package commands

import (
	"testing"
)

func TestEncode(t *testing.T) {
	cmd := Command{
		0xF00F,
		5,
		0x0ff0,
		[][]uint8{
			{0x01},
			{0x02, 0x03},
		},
	}
	data := cmd.encode()
	if data[0] != 0xf0 {
		t.Errorf("expected %x and got %x", 0xf0, data[0])
	}
}
