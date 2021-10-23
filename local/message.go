package local

import (
	"encoding/json"
)

// Basic structure of transport package.
type Message struct {
	Head HeadMessage `json:"head"`
	Body BodyMessage `json:"body"`
}

type HeadMessage struct {
	Diff    uint8  `json:"diff"`
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
func NewMessage(title, data []byte, diff uint) *Message {
	return &Message{
		Head: HeadMessage{
			Diff:  uint8(diff), // <- change WithDiff()
			Title: title,
		},
		Body: BodyMessage{
			Data: data,
		},
	}
}

// Serialize with JSON format.
func (msg *Message) Serialize() Package {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return jsonData
}
