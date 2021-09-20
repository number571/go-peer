package local

import (
	"encoding/json"
	"math"
)

// Basic structure of transport package.
type Message struct {
	Head HeadMessage `json:"head"`
	Body BodyMessage `json:"body"`
}

type HeadMessage struct {
	diff    uint8
	Title   []byte `json:"title"`
	Rand    []byte `json:"rand"`
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
func NewMessage(title, data []byte) *Message {
	return &Message{
		Head: HeadMessage{
			diff:  math.MaxUint8,
			Title: title,
		},
		Body: BodyMessage{
			Data: data,
		},
	}
}

// Append difference of message.
func (msg *Message) WithDiff(diff uint) *Message {
	msg.Head.diff = uint8(diff)
	return msg
}

// Serialize with JSON format.
func (msg *Message) Serialize() Package {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return jsonData
}
