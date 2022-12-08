package database

import (
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

type IKeyValueDB interface {
	Size(asymmetric.IPubKey) uint64
	Push(asymmetric.IPubKey, IMessage) error
	Load(asymmetric.IPubKey, uint64, uint64) ([]IMessage, error)

	modules.ICloser
}

type IMessage interface {
	IsIncoming() bool
	GetMessage() string
	GetTimestamp() string
	Bytes() []byte
}
