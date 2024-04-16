package message

import (
	"bytes"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ IMessage = &SMessage{}
)

const (
	CSeparator    = "@"
	cSeparatorLen = len(CSeparator)
)

// Basic structure of transport package.
type SMessage struct {
	FPubk string `json:"pubk"`
	FEnck string `json:"enck"`
	FSalt string `json:"salt"`
	FHash string `json:"hash"`
	FSign string `json:"sign"`
	FData []byte `json:"-"`
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
		return nil, ErrUnknownType
	}

	i := bytes.Index(recvMsg, []byte(CSeparator))
	if i == -1 {
		return nil, ErrSeparatorNotFound
	}

	if err := encoding.DeserializeJSON(recvMsg[:i], msg); err != nil {
		return nil, utils.MergeErrors(ErrDeserializeMessage, err)
	}

	switch x := pMsg.(type) {
	case []byte:
		msg.FData = x[i+cSeparatorLen:]
	case string:
		encStr := strings.TrimSpace(x[i+cSeparatorLen:])
		msg.FData = encoding.HexDecode(encStr)
		if msg.FData == nil {
			return nil, ErrDecodePayload
		}
	}

	if !msg.IsValid(psett) {
		return nil, ErrInvalidMessage
	}

	return msg, nil
}

func (p *SMessage) GetPubk() []byte {
	return encoding.HexDecode(p.FPubk)
}

func (p *SMessage) GetEnck() []byte {
	return encoding.HexDecode(p.FEnck)
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

func (p *SMessage) GetData() []byte {
	return p.FData
}

func (p *SMessage) ToBytes() []byte {
	return bytes.Join(
		[][]byte{
			encoding.SerializeJSON(p),
			p.FData,
		},
		[]byte(CSeparator),
	)
}

func (p *SMessage) ToString() string {
	return strings.Join(
		[]string{
			string(encoding.SerializeJSON(p)),
			encoding.HexEncode(p.FData),
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
		uint64(len(p.GetEnck())) != keySizeBytes,
		uint64(len(p.GetSign())) != symmetric.CAESBlockSize+keySizeBytes,
		uint64(len(p.GetPubk())) < symmetric.CAESBlockSize+keySizeBytes,
		len(p.GetHash()) != symmetric.CAESBlockSize+hashing.CSHA256Size,
		len(p.GetSalt()) != symmetric.CAESBlockSize+symmetric.CAESKeySize,
		len(p.GetData()) < symmetric.CAESBlockSize:
		return false
	default:
		return true
	}
}
