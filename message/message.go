package message

import (
	"encoding/json"

	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/settings"
)

var (
	_ IMessage = &SMessage{}
	_ iHead    = SHeadMessage{}
	_ iBody    = SBodyMessage{}
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

// IMessage

func LoadMessage(bmsg []byte) IMessage {
	var msg = new(SMessage)
	err := json.Unmarshal(bmsg, msg)
	if err != nil {
		return nil
	}
	return msg
}

func (msg *SMessage) Head() iHead {
	return msg.FHead
}

func (msg *SMessage) Body() iBody {
	return msg.FBody
}

func (msg *SMessage) Bytes() []byte {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return jsonData
}

// IHead

func (head SHeadMessage) Sender() []byte {
	return encoding.HexDecode(head.FSender)
}

func (head SHeadMessage) Session() []byte {
	return encoding.HexDecode(head.FSession)
}

func (head SHeadMessage) Salt() []byte {
	return encoding.HexDecode(head.FSalt)
}

// IBody

func (body SBodyMessage) Payload() payload.IPayload {
	return payload.LoadPayload(encoding.HexDecode(body.FPayload))
}

func (body SBodyMessage) Hash() []byte {
	return encoding.HexDecode(body.FHash)
}

func (body SBodyMessage) Sign() []byte {
	return encoding.HexDecode(body.FSign)
}

func (body SBodyMessage) Proof() uint64 {
	bProof := encoding.HexDecode(body.FProof)
	if len(bProof) != settings.CSizeUint64 {
		return 0
	}
	res := [settings.CSizeUint64]byte{}
	copy(res[:], bProof)
	return encoding.BytesToUint64(res)
}
