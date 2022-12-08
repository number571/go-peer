package database

import (
	"bytes"
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	// TODO: need auto time
	fIsIncoming bool
	fMessage    string
}

func NewMessage(isIncoming bool, message string) IMessage {
	return &sMessage{isIncoming, message}
}

func LoadMessage(msgBytes []byte) IMessage {
	if len(msgBytes) == 0 {
		return nil
	}
	isIncoming := false
	if msgBytes[0] == 1 {
		isIncoming = true
	}
	return &sMessage{isIncoming, string(msgBytes[1:])}
}

func (msg *sMessage) IsIncoming() bool {
	return msg.fIsIncoming
}

func (msg *sMessage) GetMessage() string {
	return msg.fMessage
}

func (msg *sMessage) Bytes() []byte {
	firstByte := byte(0)
	if msg.fIsIncoming {
		firstByte = byte(1)
	}
	return bytes.Join(
		[][]byte{
			{firstByte},
			[]byte(msg.fMessage),
		},
		[]byte{},
	)
}
