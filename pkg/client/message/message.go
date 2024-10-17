package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/encoding"
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
	var recvMsg []byte
	switch x := pMsg.(type) {
	case []byte:
		recvMsg = x
	case string:
		recvMsg = encoding.HexDecode(x)
	default:
		return nil, ErrUnknownMessageType
	}
	return loadMessage(psett, recvMsg)
}

func (p *sMessage) GetEnck() []byte {
	return p.fEnck
}

func (p *sMessage) GetEncd() []byte {
	return p.fEncd
}

func (p *sMessage) ToBytes() []byte {
	return bytes.Join(
		[][]byte{p.fEnck, p.fEncd},
		[]byte{},
	)
}

func (p *sMessage) ToString() string {
	return encoding.HexEncode(p.ToBytes())
}

func loadMessage(pSett ISettings, pBytes []byte) (IMessage, error) {
	keySize := pSett.GetEncKeySizeBytes()
	msgSize := pSett.GetMessageSizeBytes()
	if keySize >= msgSize {
		panic("keySize >= msgSize")
	}
	if uint64(len(pBytes)) != msgSize {
		return nil, ErrLoadMessageBytes
	}
	return NewMessage(pBytes[:keySize], pBytes[keySize:]), nil
}
