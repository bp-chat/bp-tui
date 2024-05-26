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
	expected := []uint8{0xf0, 0x0f, 0x0a, 0x05, 0x0a, 0x0f, 0xf0, 0x0a, 0x01, 0x0a, 0x02, 0x03}
	data, err := cmd.encode()
	if err != nil {
		t.Errorf("error while encoding the data")
		t.Error(err)
		return
	}
	if len(data) != len(expected) {
		t.Fatalf("Different len between data %v and array %v", data, expected)
	}
	for index, element := range expected {

		if data[index] != element {
			t.Errorf("expected %x and got %x", element, data[index])
		}
	}
}

func TestDecode(t *testing.T) {
	source := []uint8{0xf0, 0x0f, 0x0a, 0x05, 0x0a, 0x0f, 0xf0, 0x0a, 0x01, 0x0a, 0x02, 0x03}
	actual, err := decode(source, 2)
	if err != nil {
		t.Fatalf("Could not decode data %v", err)
	}
	if actual.version != 0xf00f {
		t.Errorf("Could not parse id, expected %d and got %d", 0xf00f, actual.id)
	}
	if actual.sync_id != 5 {
		t.Errorf("Could not parse sync id, expected %d and got %d", 0x05, actual.sync_id)
	}
	if actual.id != 0x0ff0 {
		t.Errorf("Could not parse version, expected %d and got %d", 0x0ff0, actual.version)
	}
	if actual.args[0][0] != 0x01 {
		t.Errorf("Could not parse first argument, expected %d and got %d", 0x01, actual.args[0][0])
	}
	if actual.args[1][0] != 0x02 || actual.args[1][1] != 0x03 {
		t.Errorf("Could not parse second argument, expected %d and got %d", 0x02, actual.args[1][0])
		t.Errorf("Could not parse second argument, expected %d and got %d", 0x03, actual.args[1][1])
	}
}
