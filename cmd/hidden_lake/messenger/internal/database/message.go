package database

import (
	"bytes"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
)

const (
	cIsIncomingSize = 1
	cTimestampSize  = encoding.CSizeUint64
	cBlockUIDSize   = 16
)

var (
	_ IMessage = &sMessage{}
)

type sMessage struct {
	fIsIncoming bool
	fTimestamp  uint64
	fBlockUID   [cBlockUIDSize]byte
	fMessage    []byte
}

func NewMessage(pIsIncoming bool, pMessage []byte) IMessage {
	blockUID := [cBlockUIDSize]byte{}
	copy(blockUID[:], random.NewStdPRNG().GetBytes(cBlockUIDSize))

	return &sMessage{
		fIsIncoming: pIsIncoming,
		fBlockUID:   blockUID,
		fTimestamp:  uint64(time.Now().Unix()),
		fMessage:    pMessage,
	}
}

func LoadMessage(pMsgBytes []byte) IMessage {
	if len(pMsgBytes) < (cIsIncomingSize + cTimestampSize + cBlockUIDSize) {
		return nil
	}

	isIncoming := false
	if pMsgBytes[0] == 1 {
		isIncoming = true
	}

	blockTimestamp := [cTimestampSize]byte{}
	copy(blockTimestamp[:], pMsgBytes[cIsIncomingSize:cIsIncomingSize+cTimestampSize])

	blockUID := [cBlockUIDSize]byte{}
	copy(blockUID[:], pMsgBytes[cIsIncomingSize:cIsIncomingSize+cTimestampSize:cIsIncomingSize+cTimestampSize+cBlockUIDSize])

	return &sMessage{
		fIsIncoming: isIncoming,
		fBlockUID:   blockUID,
		fTimestamp:  encoding.BytesToUint64(blockTimestamp),
		fMessage:    pMsgBytes[cIsIncomingSize+cTimestampSize+cBlockUIDSize:],
	}
}

func (p *sMessage) IsIncoming() bool {
	return p.fIsIncoming
}

func (p *sMessage) GetBlockUID() [cBlockUIDSize]byte {
	return p.fBlockUID
}

func (p *sMessage) GetMessage() []byte {
	return p.fMessage
}

func (p *sMessage) GetTimestamp() string {
	return time.Unix(int64(p.fTimestamp), 0).Format("2006-01-02T15:04:05")
}

func (p *sMessage) ToBytes() []byte {
	firstByte := byte(0)
	if p.fIsIncoming {
		firstByte = byte(1)
	}
	blockTimestamp := encoding.Uint64ToBytes(p.fTimestamp)
	return bytes.Join(
		[][]byte{
			{firstByte},
			blockTimestamp[:],
			[]byte(p.fBlockUID[:]),
			[]byte(p.fMessage),
		},
		[]byte{},
	)
}
