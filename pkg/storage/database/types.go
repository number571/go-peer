package database

import (
	"github.com/number571/go-peer/pkg/types"
)

type IKVDatabase interface {
	types.ICloser
	GetSettings() ISettings

	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Del([]byte) error
}

type ISettings interface {
	GetPath() string
	GetWorkSize() uint64
	GetPassword() string
}
