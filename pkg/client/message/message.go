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
	cSeparator    = "==="
	cSeparatorLen = len(cSeparator)
)

var (
	_ IMessage = &SMessage{}
	_ IHead    = SHeadMessage{}
	_ IBody    = SBodyMessage{}
)

// Basic structure of transport package.
type SMessage struct {
	FHead SHeadMessage `json:"head"`
	FBody SBodyMessage `json:"body"`
}

type SHeadMessage struct {
	FSalt    string `json:"salt"`
	FSession string `json:"session"`
	FSender  string `json:"sender"`
}

type SBodyMessage struct {
	FSign    string `json:"sign"`
	FHash    string `json:"hash"`
	FProof   string `json:"proof"`
	FPayload []byte `json:"-"`
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
	}

	i := bytes.Index(recvMsg, []byte(cSeparator))
	if i == -1 {
		return nil
	}

	if err := json.Unmarshal(recvMsg[:i], msg); err != nil {
		return nil
	}

	switch x := pMsg.(type) {
	case []byte:
		msg.FBody.FPayload = x[i+cSeparatorLen:]
	case string:
		encStr := strings.TrimSpace(x[i+cSeparatorLen:])
		msg.FBody.FPayload = encoding.HexDecode(encStr)
		if msg.FBody.FPayload == nil {
			return nil
		}
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
			[]byte(cSeparator),
			p.FBody.FPayload,
		},
		[]byte{},
	)
}

func (p *SMessage) ToString() string {
	return strings.Join(
		[]string{
			string(encoding.Serialize(p, false)),
			cSeparator,
			encoding.HexEncode(p.FBody.FPayload),
		},
		"",
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

func (p SBodyMessage) GetPayload() payload.IPayload {
	return payload.LoadPayload(p.FPayload)
}

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
