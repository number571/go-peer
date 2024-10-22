package message

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
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
	if len(pEnck) != asymmetric.CKEMCiphertextSize {
		panic(`len(pEnck) != asymmetric.CKEMCiphertextSize`)
	}
	return &sMessage{
		fEnck: pEnck,
		fEncd: pEncd,
	}
}

// Message can be created only with client module.
func LoadMessage(pSize uint64, pMsg interface{}) (IMessage, error) {
	kSize := uint64(asymmetric.CKEMCiphertextSize)
	if kSize >= pSize {
		return nil, ErrSizeMessageBytes
	}
	var recvMsg []byte
	switch x := pMsg.(type) {
	case []byte:
		recvMsg = x
	case string:
		recvMsg = encoding.HexDecode(x)
	default:
		return nil, ErrUnknownMessageType
	}
	if uint64(len(recvMsg)) != pSize {
		return nil, ErrLoadMessageBytes
	}
	return NewMessage(recvMsg[:kSize], recvMsg[kSize:]), nil
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
