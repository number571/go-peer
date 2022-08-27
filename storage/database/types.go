package database

import (
	"github.com/number571/go-peer/storage"
)

type IKeyValueDB interface {
	storage.IKeyValueStorage
	Close() error

	Settings() ISettings
	Iter([]byte) iIterator
}

type ISettings interface {
	GetPath() string
	GetHashing() bool
	GetCipherKey() []byte
}

type iIterator interface {
	Key() []byte
	Value() []byte

	Next() bool
	Close()
}
