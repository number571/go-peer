package database

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"
)

type IWrapperDB interface {
	types.ICloser

	Get() IKeyValueDB
	Set(IKeyValueDB) IWrapperDB
}

type IKeyValueDB interface {
	types.ICloser

	Size(IRelation) uint64
	Push(IRelation, IMessage) error
	Load(IRelation, uint64, uint64) ([]IMessage, error)
}

type IRelation interface {
	IAm() asymmetric.IPubKey
	Friend() asymmetric.IPubKey
}

type IMessage interface {
	IsIncoming() bool
	GetSHA256UID() string
	GetMessage() string
	GetTimestamp() string
	ToBytes() []byte
}
