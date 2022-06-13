package message

import (
	"bytes"
	"encoding/json"

	"github.com/number571/go-peer/encoding"
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
	FData  []byte `json:"data"`
	FHash  []byte `json:"hash"`
	FSign  []byte `json:"sign"`
	FProof uint64 `json:"proof"`
}

// IMessage

func NewMessage(title, data []byte) IMessage {
	return &SMessage{
		FHead: SHeadMessage{},
		FBody: SBodyMessage{
			FData: bytes.Join([][]byte{
				encoding.Uint64ToBytes(uint64(len(title))),
				title,
				data,
			}, []byte{}),
		},
	}
}

func (msg *SMessage) Head() iHead {
	return msg.FHead
}

func (msg *SMessage) Body() iBody {
	return msg.FBody
}

func (msg *SMessage) ToPackage() IPackage {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return LoadPackage(jsonData)
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

func (body SBodyMessage) Data() []byte {
	return body.FData
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
