package local

import "encoding/json"

// Basic structure of transport package.
type Message struct {
	Head HeadMessage `json:"head"`
	Body BodyMessage `json:"body"`
}

type HeadMessage struct {
	Title   []byte `json:"title"`
	Rand    []byte `json:"rand"`
	Sender  []byte `json:"sender"`
	Session []byte `json:"session"`
	Diff    uint8  `json:"diff"`
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
			Title: title,
		},
		Body: BodyMessage{
			Data: data,
		},
	}
}

// Append difference of message.
func (msg *Message) WithDiff(diff uint) *Message {
	msg.Head.Diff = uint8(diff)
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
