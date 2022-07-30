package message

import (
	"encoding/json"

	"github.com/number571/go-peer/local/payload"
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
	FSender  []byte `json:"sender"`
	FSession []byte `json:"session"`
	FSalt    []byte `json:"salt"`
}

type SBodyMessage struct {
	FPayload []byte `json:"payload"`
	FHash    []byte `json:"hash"`
	FSign    []byte `json:"sign"`
	FProof   uint64 `json:"proof"`
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
	return head.FSender
}

func (head SHeadMessage) Session() []byte {
	return head.FSession
}

func (head SHeadMessage) Salt() []byte {
	return head.FSalt
}

// IBody

func (body SBodyMessage) Payload() payload.IPayload {
	return payload.LoadPayload(body.FPayload)
}

func (body SBodyMessage) Hash() []byte {
	return body.FHash
}

func (body SBodyMessage) Sign() []byte {
	return body.FSign
}

func (body SBodyMessage) Proof() uint64 {
	return body.FProof
}
