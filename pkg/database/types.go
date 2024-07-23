package database

import (
	"github.com/number571/go-peer/pkg/types"
)

type IKVDatabase interface {
	types.ICloser

	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Del([]byte) error
}
