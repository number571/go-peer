package database

import (
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/types"
)

type ISettings interface {
	GetPath() string
	GetHashing() bool
	GetCipherKey() []byte
}

type IKVDatabase interface {
	storage.IKVStorage
	types.ICloser

	GetSettings() ISettings
}
