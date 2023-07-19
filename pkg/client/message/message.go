package message

import (
	"encoding/json"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
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
	FPayload string `json:"payload"`
	FSign    string `json:"sign"`
	FHash    string `json:"hash"`
	FProof   string `json:"proof"`
}

// Message can be created only with client module.
func LoadMessage(psett ISettings, pMsg []byte) IMessage {
	msg := new(SMessage)
	if err := json.Unmarshal(pMsg, msg); err != nil {
		return nil
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
	jsonData, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return jsonData
}

func (p *SMessage) ToString() string {
	return string(p.ToBytes())
}

func (p *SMessage) IsValid(psett ISettings) bool {
	if uint64(len(p.ToBytes())) > psett.GetMessageSizeBytes() {
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
	return payload.LoadPayload(encoding.HexDecode(p.FPayload))
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
