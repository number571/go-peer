package database

import (
	"bytes"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fIsIncoming bool
	fSHA256UID  []byte
	fTimestamp  string
	fMessage    string
}

func NewMessage(isIncoming bool, message string, hashUID []byte) IMessage {
	t := time.Now()
	return &sMessage{
		fIsIncoming: isIncoming,
		fSHA256UID:  hashUID,
		// 2022-12-08T02:43:24
		fTimestamp: fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second()),
		fMessage: message,
	}
}

func LoadMessage(msgBytes []byte) IMessage {
	// 1(isIncoming) + 32(sha256) + 19(time)
	if len(msgBytes) < 1+hashing.CSHA256Size+19 {
		return nil
	}
	isIncoming := false
	if msgBytes[0] == 1 {
		isIncoming = true
	}
	return &sMessage{
		fIsIncoming: isIncoming,                                                         // 1 byte
		fSHA256UID:  msgBytes[1 : hashing.CSHA256Size+1],                                // 32 bytes
		fTimestamp:  string(msgBytes[hashing.CSHA256Size+1 : hashing.CSHA256Size+1+19]), // 19 bytes
		fMessage:    string(msgBytes[hashing.CSHA256Size+1+19:]),                        // n bytes
	}
}

func (msg *sMessage) IsIncoming() bool {
	return msg.fIsIncoming
}

func (msg *sMessage) GetSHA256UID() string {
	return string(msg.fSHA256UID)
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
			[]byte(msg.fSHA256UID),
			[]byte(msg.fTimestamp),
			[]byte(msg.fMessage),
		},
		[]byte{},
	)
}
