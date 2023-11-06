package message

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	CSeparator    = "==="
	cSeparatorLen = len(CSeparator)
)

var (
	_ IMessage = &SMessage{}
	_ IHead    = SHeadMessage{}
	_ IBody    = SBodyMessage{}
)

// Basic structure of transport package.
type SMessage struct {
	FHead    SHeadMessage `json:"head"`
	FBody    SBodyMessage `json:"body"`
	FPayload []byte       `json:"-"`
}

type SHeadMessage struct {
	FSalt    string `json:"salt"`
	FSession string `json:"session"`
	FSender  string `json:"sender"`
}

type SBodyMessage struct {
	FSign  string `json:"sign"`
	FHash  string `json:"hash"`
	FProof string `json:"proof"`
}

// Message can be created only with client module.
func LoadMessage(psett ISettings, pMsg interface{}) IMessage {
	msg := new(SMessage)

	var recvMsg []byte
	switch x := pMsg.(type) {
	case []byte:
		recvMsg = x
	case string:
		recvMsg = []byte(x)
	default:
		return nil
	}

	i := bytes.Index(recvMsg, []byte(CSeparator))
	if i == -1 {
		return nil
	}

	if err := json.Unmarshal(recvMsg[:i], msg); err != nil {
		return nil
	}

	switch x := pMsg.(type) {
	case []byte:
		msg.FPayload = x[i+cSeparatorLen:]
	case string:
		encStr := strings.TrimSpace(x[i+cSeparatorLen:])
		msg.FPayload = encoding.HexDecode(encStr)
		if msg.FPayload == nil {
			return nil
		}
	default:
		panic("got unknown type")
	}

	if !msg.IsValid(psett) {
		return nil
	}

	return msg
}

func (p *SMessage) GetHead() IHead {
	return p.FHead
}

func (p *SMessage) GetBody() IBody {
	return p.FBody
}

func (p *SMessage) ToBytes() []byte {
	return bytes.Join(
		[][]byte{
			encoding.Serialize(p, false),
			p.FPayload,
		},
		[]byte(CSeparator),
	)
}

func (p *SMessage) ToString() string {
	return strings.Join(
		[]string{
			string(encoding.Serialize(p, false)),
			encoding.HexEncode(p.FPayload),
		},
		CSeparator,
	)
}

func (p *SMessage) IsValid(psett ISettings) bool {
	if uint64(len(p.ToBytes())) != psett.GetMessageSizeBytes() {
		return false
	}
	if len(p.GetBody().GetHash()) != hashing.CSHA256Size {
		return false
	}
	puzzle := puzzle.NewPoWPuzzle(psett.GetWorkSizeBits())
	return puzzle.VerifyBytes(p.GetBody().GetHash(), p.GetBody().GetProof())
}

func (p *SMessage) GetPayload() payload.IPayload {
	return payload.LoadPayload(p.FPayload)
}

// IHead

func (p SHeadMessage) GetSender() []byte {
	return encoding.HexDecode(p.FSender)
}

func (p SHeadMessage) GetSession() []byte {
	return encoding.HexDecode(p.FSession)
}

func (p SHeadMessage) GetSalt() []byte {
	return encoding.HexDecode(p.FSalt)
}

// IBody

func (p SBodyMessage) GetHash() []byte {
	return encoding.HexDecode(p.FHash)
}

func (p SBodyMessage) GetSign() []byte {
	return encoding.HexDecode(p.FSign)
}

func (p SBodyMessage) GetProof() uint64 {
	bProof := encoding.HexDecode(p.FProof)
	if len(bProof) != encoding.CSizeUint64 {
		return 0
	}
	res := [encoding.CSizeUint64]byte{}
	copy(res[:], bProof)
	return encoding.BytesToUint64(res)
}
