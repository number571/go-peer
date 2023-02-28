package database

import (
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/types"
)

type ISettings interface {
	GetPath() string
	GetHashing() bool
	GetSaltKey() []byte
	GetCipherKey() []byte
}

type IKeyValueDB interface {
	storage.IKeyValueStorage
	types.ICloser

	GetSettings() ISettings
	GetIterator([]byte) IIterator
}

type IIterator interface {
	types.ICloser
	Next() bool

	GetKey() []byte
	GetValue() []byte
}
