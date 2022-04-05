package local

import (
	"bytes"
	"encoding/json"

	"github.com/number571/go-peer/encoding"
)

var (
	_ IMessage = &sMessage{}
	_ iHead    = sHeadMessage{}
	_ iBody    = sBodyMessage{}
)

// Basic structure of transport package.
type sMessage struct {
	FHead sHeadMessage `json:"head"`
	FBody sBodyMessage `json:"body"`
}

type sHeadMessage struct {
	FSender  []byte `json:"sender"`
	FSession []byte `json:"session"`
	FSalt    []byte `json:"salt"`
}

type sBodyMessage struct {
	FData  []byte `json:"data"`
	FHash  []byte `json:"hash"`
	FSign  []byte `json:"sign"`
	FProof uint64 `json:"proof"`
}

// IMessage

func NewMessage(title, data []byte) IMessage {
	return &sMessage{
		FHead: sHeadMessage{},
		FBody: sBodyMessage{
			FData: bytes.Join([][]byte{
				encoding.Uint64ToBytes(uint64(len(title))),
				title,
				data,
			}, []byte{}),
		},
	}
}

func (msg *sMessage) Head() iHead {
	return msg.FHead
}

func (msg *sMessage) Body() iBody {
	return msg.FBody
}

func (msg *sMessage) export() []byte {
	const (
		SizeUint64 = 8 // bytes
	)

	if len(msg.FBody.FData) < SizeUint64 {
		return nil
	}

	mustLen := encoding.BytesToUint64(msg.FBody.FData[:SizeUint64])
	allData := msg.FBody.FData[SizeUint64:]
	if mustLen > uint64(len(allData)) {
		return nil
	}

	msg.FBody.FData = allData[mustLen:]
	return allData[:mustLen]
}

func (msg *sMessage) ToPackage() IPackage {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return LoadPackage(jsonData)
}

// IHead

func (head sHeadMessage) Sender() []byte {
	return head.FSender
}

func (head sHeadMessage) Session() []byte {
	return head.FSession
}

func (head sHeadMessage) Salt() []byte {
	return head.FSalt
}

// IBody

func (body sBodyMessage) Data() []byte {
	return body.FData
}

func (body sBodyMessage) Hash() []byte {
	return body.FHash
}

func (body sBodyMessage) Sign() []byte {
	return body.FSign
}

func (body sBodyMessage) Proof() uint64 {
	return body.FProof
}
