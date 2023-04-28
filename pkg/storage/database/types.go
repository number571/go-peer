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
}
