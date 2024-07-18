package commands

type UserName = [16]byte

type Header struct {
	Version uint16
	SyncId  uint8
	Id      CommandCodes
}

const headerSize = 5

const version uint16 = 0

type CommandCodes uint16

const (
	Undefined CommandCodes = 0
	CNN                    = 1 //Connect
	MSG                    = 2 //Message
	RKS                    = 3 //Register Keys
	BKS                    = 4 //Broadcast Keys, this command is temporary and must be removed
)

func NewHeader(id CommandCodes, sync uint8) Header {
	return Header{
		version,
		sync,
		id,
	}
}
