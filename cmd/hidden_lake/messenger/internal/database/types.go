package database

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"
)

type IWrapperDB interface {
	Get() IKeyValueDB
	Update(IKeyValueDB) error

	types.ICloser
}

type IKeyValueDB interface {
	Size(IRelation) uint64
	Push(IRelation, IMessage) error
	Load(IRelation, uint64, uint64) ([]IMessage, error)

	types.ICloser
}

type IRelation interface {
	IAm() asymmetric.IPubKey
	Friend() asymmetric.IPubKey
}

type IMessage interface {
	IsIncoming() bool
	GetMessage() string
	GetTimestamp() string
	Bytes() []byte
}
