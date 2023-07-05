package database

import (
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/types"
)

type IKVDatabase interface {
	storage.IKVStorage
	types.ICloser
}
