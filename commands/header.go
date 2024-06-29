package commands

type Header struct {
	Version uint16
	SyncId  uint8
	Id      CommandCode
}

const version uint16 = 0

type CommandCode uint16

const (
	Undefined CommandCode = 0
	CNN                   = 1
	MSG                   = 2
	RKS                   = 3
)

func NewHeader(id CommandCode, sync uint8) Header {
	return Header{
		version,
		sync,
		id,
	}
}
