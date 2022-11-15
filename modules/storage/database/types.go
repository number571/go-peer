package database

import (
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/storage"
)

type IKeyValueDB interface {
	storage.IKeyValueStorage
	modules.ICloser

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
