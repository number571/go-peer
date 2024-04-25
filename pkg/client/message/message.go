package message

import (
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

var (
	_ IMessage = &sMessage{}
)

// Basic structure of transport package.
type sMessage struct {
	fEnck []byte
	fEncd []byte
}

func NewMessage(pEnck, pEncd []byte) IMessage {
	return &sMessage{
		fEnck: pEnck,
		fEncd: pEncd,
	}
}

// Message can be created only with client module.
func LoadMessage(psett ISettings, pMsg interface{}) (IMessage, error) {
	msg := new(sMessage)

	var recvMsg []byte
	switch x := pMsg.(type) {
	case []byte:
		recvMsg = x
	case string:
		recvMsg = encoding.HexDecode(x)
	default:
		return nil, ErrUnknownType
	}

	wrapSlice, err := joiner.LoadBytesJoiner(recvMsg)
	if err != nil || len(wrapSlice) != 2 {
		return nil, ErrLoadBytesJoiner
	}

	msg.fEnck = wrapSlice[0]
	msg.fEncd = wrapSlice[1]

	if !msg.IsValid(psett) {
		return nil, ErrInvalidMessage
	}

	return msg, nil
}

func (p *sMessage) GetEnck() []byte {
	return p.fEnck
}

func (p *sMessage) GetEncd() []byte {
	return p.fEncd
}

func (p *sMessage) ToBytes() []byte {
	return joiner.NewBytesJoiner([][]byte{p.fEnck, p.fEncd})
}

func (p *sMessage) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}

func (p *sMessage) IsValid(psett ISettings) bool {
	msgSizeBytes := psett.GetMessageSizeBytes()
	keySizeBytes := psett.GetKeySizeBits() / 8
	switch {
	case
		uint64(len(p.ToBytes())) != msgSizeBytes,
		uint64(len(p.GetEnck())) != keySizeBytes:
		return false
	default:
		return true
	}
}
