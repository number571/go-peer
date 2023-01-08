package database

import (
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/types"
)

type IKeyValueDB interface {
	storage.IKeyValueStorage
	types.ICloser

	Settings() ISettings
	Iter([]byte) IIterator
}

type ISettings interface {
	GetHashing() bool
	GetCipherKey() []byte
}

type IIterator interface {
	Key() []byte
	Value() []byte

	Next() bool
	Close()
}
