package database

import (
	"bytes"
	"fmt"
	"time"
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fIsIncoming bool
	fMessage    string
	fTimestamp  string
}

func NewMessage(isIncoming bool, message string) IMessage {
	t := time.Now()
	return &sMessage{
		fIsIncoming: isIncoming,
		fMessage:    message,
		// 2022-12-08T02:43:24
		fTimestamp: fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second()),
	}
}

func LoadMessage(msgBytes []byte) IMessage {
	if len(msgBytes) < 20 { // 1(isIncoming)+19(time)
		return nil
	}
	isIncoming := false
	if msgBytes[0] == 1 {
		isIncoming = true
	}
	return &sMessage{
		fIsIncoming: isIncoming,
		fTimestamp:  string(msgBytes[1:20]),
		fMessage:    string(msgBytes[20:]),
	}
}

func (msg *sMessage) IsIncoming() bool {
	return msg.fIsIncoming
}

func (msg *sMessage) GetMessage() string {
	return msg.fMessage
}

func (msg *sMessage) GetTimestamp() string {
	return msg.fTimestamp
}

func (msg *sMessage) Bytes() []byte {
	firstByte := byte(0)
	if msg.fIsIncoming {
		firstByte = byte(1)
	}
	return bytes.Join(
		[][]byte{
			{firstByte},
			[]byte(msg.fTimestamp),
			[]byte(msg.fMessage),
		},
		[]byte{},
	)
}
