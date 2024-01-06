package message

import (
	"bytes"
	"errors"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IMessage = &SMessage{}
)

const (
	CSeparator    = "@"
	cSeparatorLen = len(CSeparator)
)

var (
	_ IMessage = &SMessage{}
)

// Basic structure of transport package.
type SMessage struct {
	FPubKey  string `json:"pubk"`
	FEncKey  string `json:"enck"`
	FSalt    string `json:"salt"`
	FHash    string `json:"hash"`
	FSign    string `json:"sign"`
	FPayload []byte `json:"-"`
}

// Message can be created only with client module.
func LoadMessage(psett ISettings, pMsg interface{}) (IMessage, error) {
	msg := new(SMessage)

	var recvMsg []byte
	switch x := pMsg.(type) {
	case []byte:
		recvMsg = x
	case string:
		recvMsg = []byte(x)
	default:
		return nil, errors.New("unknown type of message")
	}

	i := bytes.Index(recvMsg, []byte(CSeparator))
	if i == -1 {
		return nil, errors.New("undefined separator")
	}

	if err := encoding.DeserializeJSON(recvMsg[:i], msg); err != nil {
		return nil, errors.New("unmarshal message")
	}

	switch x := pMsg.(type) {
	case []byte:
		msg.FPayload = x[i+cSeparatorLen:]
	case string:
		encStr := strings.TrimSpace(x[i+cSeparatorLen:])
		msg.FPayload = encoding.HexDecode(encStr)
		if msg.FPayload == nil {
			return nil, errors.New("hex decode payload")
		}
	}

	if !msg.IsValid(psett) {
		return nil, errors.New("message is invalid")
	}

	return msg, nil
}

func (p *SMessage) GetPubKey() []byte {
	return encoding.HexDecode(p.FPubKey)
}

func (p *SMessage) GetEncKey() []byte {
	return encoding.HexDecode(p.FEncKey)
}

func (p *SMessage) GetSalt() []byte {
	return encoding.HexDecode(p.FSalt)
}

func (p *SMessage) GetHash() []byte {
	return encoding.HexDecode(p.FHash)
}

func (p *SMessage) GetSign() []byte {
	return encoding.HexDecode(p.FSign)
}

func (p *SMessage) GetPayload() []byte {
	return p.FPayload
}

func (p *SMessage) ToBytes() []byte {
	return bytes.Join(
		[][]byte{
			encoding.SerializeJSON(p),
			p.FPayload,
		},
		[]byte(CSeparator),
	)
}

func (p *SMessage) ToString() string {
	return strings.Join(
		[]string{
			string(encoding.SerializeJSON(p)),
			encoding.HexEncode(p.FPayload),
		},
		CSeparator,
	)
}

func (p *SMessage) IsValid(psett ISettings) bool {
	msgSizeBytes := psett.GetMessageSizeBytes()
	keySizeBytes := psett.GetKeySizeBits() / 8
	switch {
	case
		uint64(len(p.ToBytes())) != msgSizeBytes,
		uint64(len(p.GetEncKey())) != keySizeBytes,
		uint64(len(p.GetSign())) != symmetric.CAESBlockSize+keySizeBytes,
		uint64(len(p.GetPubKey())) < symmetric.CAESBlockSize+keySizeBytes,
		len(p.GetHash()) != symmetric.CAESBlockSize+hashing.CSHA256Size,
		len(p.GetSalt()) != symmetric.CAESBlockSize+symmetric.CAESKeySize,
		len(p.GetPayload()) < symmetric.CAESBlockSize:
		return false
	default:
		return true
	}
}
