package database

import (
	"github.com/number571/go-peer/local"
)

type IKeyValueDB interface {
	Size([]byte) uint64
	Push([]byte, local.IMessage) error
	Load([]byte, uint64) local.IMessage

	Close() error
	Clean() error
}
