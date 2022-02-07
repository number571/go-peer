package local

import (
	"bytes"
	"encoding/json"

	"github.com/number571/go-peer/encoding"
)

// Basic structure of transport package.
type MessageT struct {
	Head HeadMessage `json:"head"`
	Body BodyMessage `json:"body"`
}

type HeadMessage struct {
	Sender  []byte `json:"sender"`
	Session []byte `json:"session"`
}

type BodyMessage struct {
	Data []byte `json:"data"`
	Hash []byte `json:"hash"`
	Sign []byte `json:"sign"`
	Npow uint64 `json:"npow"`
}

// Create message with title and data.
func NewMessage(title, data []byte) Message {
	return &MessageT{
		Body: BodyMessage{
			Data: bytes.Join([][]byte{
				encoding.Uint64ToBytes(uint64(len(title))),
				title,
				data,
			}, []byte{}),
		},
	}
}

func (msg *MessageT) Export() ([]byte, []byte) {
	const (
		SizeUint64 = 8 // bytes
	)

	if len(msg.Body.Data) < SizeUint64 {
		return nil, nil
	}

	mustLen := encoding.BytesToUint64(msg.Body.Data[:SizeUint64])
	allData := msg.Body.Data[SizeUint64:]
	if mustLen > uint64(len(allData)) {
		return nil, nil
	}

	return allData[:mustLen], allData[mustLen:]
}

// Serialize with JSON format.
func (msg *MessageT) Serialize() Package {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return jsonData
}
