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
func LoadMessage(bmsg []byte, params IParams) IMessage {
	msg := new(SMessage)
	if err := json.Unmarshal(bmsg, msg); err != nil {
		return nil
	}
	if !msg.IsValid(params) {
		return nil
	}
	return msg
}

func (msg *SMessage) GetHead() IHead {
	return msg.FHead
}

func (msg *SMessage) GetBody() IBody {
	return msg.FBody
}

func (msg *SMessage) ToBytes() []byte {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return jsonData
}

func (msg *SMessage) IsValid(params IParams) bool {
	if uint64(len(msg.ToBytes())) > params.GetMessageSize() {
		return false
	}
	if len(msg.GetBody().GetHash()) != hashing.CSHA256Size {
		return false
	}
	puzzle := puzzle.NewPoWPuzzle(params.GetWorkSize())
	return puzzle.VerifyBytes(msg.GetBody().GetHash(), msg.GetBody().GetProof())
}

// IHead

func (head SHeadMessage) GetSender() []byte {
	return encoding.HexDecode(head.FSender)
}

func (head SHeadMessage) GetSession() []byte {
	return encoding.HexDecode(head.FSession)
}

func (head SHeadMessage) GetSalt() []byte {
	return encoding.HexDecode(head.FSalt)
}

// IBody

func (body SBodyMessage) GetPayload() payload.IPayload {
	return payload.LoadPayload(encoding.HexDecode(body.FPayload))
}

func (body SBodyMessage) GetHash() []byte {
	return encoding.HexDecode(body.FHash)
}

func (body SBodyMessage) GetSign() []byte {
	return encoding.HexDecode(body.FSign)
}

func (body SBodyMessage) GetProof() uint64 {
	bProof := encoding.HexDecode(body.FProof)
	if len(bProof) != encoding.CSizeUint64 {
		return 0
	}
	res := [encoding.CSizeUint64]byte{}
	copy(res[:], bProof)
	return encoding.BytesToUint64(res)
}
