package database

import (
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/types"
)

type IKVDatabase interface {
	storage.IKVStorage
	types.ICloser

	// NewBatch() IKVBatch
}

// type IKVBatch interface {
// 	Set([]byte, []byte) error
// 	Del([]byte) error
// 	Do() error
// }
