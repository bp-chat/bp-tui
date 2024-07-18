package commands

type UserName = [16]byte

type Header struct {
	Version uint16
	SyncId  uint8
	Id      CommandCode
}

const headerSize = 5

const version uint16 = 0

type CommandCode uint16

const (
	Undefined CommandCode = 0
	CNN                   = 1 //Connect
	MSG                   = 2 //Message
	RKS                   = 3 //Register Keys
	BKS                   = 4 //Broadcast Keys, this command is temporary and must be removed
)

func NewHeader(id CommandCode, sync uint8) Header {
	return Header{
		version,
		sync,
		id,
	}
}
