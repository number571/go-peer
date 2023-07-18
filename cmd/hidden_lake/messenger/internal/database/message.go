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
	fMessage    []byte
}

func NewMessage(pIsIncoming bool, pMessage []byte, pHashUID []byte) IMessage {
	t := time.Now()
	return &sMessage{
		fIsIncoming: pIsIncoming,
		fSHA256UID:  pHashUID,
		// 2022-12-08T02:43:24
		fTimestamp: fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second()),
		fMessage: pMessage,
	}
}

func LoadMessage(pMsgBytes []byte) IMessage {
	// 1(isIncoming) + 32(sha256) + 19(time)
	if len(pMsgBytes) < 1+hashing.CSHA256Size+19 {
		return nil
	}
	isIncoming := false
	if pMsgBytes[0] == 1 {
		isIncoming = true
	}
	return &sMessage{
		fIsIncoming: isIncoming,                                                          // 1 byte
		fSHA256UID:  pMsgBytes[1 : hashing.CSHA256Size+1],                                // 32 bytes
		fTimestamp:  string(pMsgBytes[hashing.CSHA256Size+1 : hashing.CSHA256Size+1+19]), // 19 bytes
		fMessage:    pMsgBytes[hashing.CSHA256Size+1+19:],                                // n bytes
	}
}

func (p *sMessage) IsIncoming() bool {
	return p.fIsIncoming
}

func (p *sMessage) GetSHA256UID() string {
	return string(p.fSHA256UID)
}

func (p *sMessage) GetMessage() []byte {
	return p.fMessage
}

func (p *sMessage) GetTimestamp() string {
	return p.fTimestamp
}

func (p *sMessage) ToBytes() []byte {
	firstByte := byte(0)
	if p.fIsIncoming {
		firstByte = byte(1)
	}
	return bytes.Join(
		[][]byte{
			{firstByte},
			[]byte(p.fSHA256UID),
			[]byte(p.fTimestamp),
			[]byte(p.fMessage),
		},
		[]byte{},
	)
}
