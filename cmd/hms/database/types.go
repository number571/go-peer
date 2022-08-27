package database

import (
	"github.com/number571/go-peer/modules/message"
)

type IKeyValueDB interface {
	Size([]byte) (uint64, error)
	Push([]byte, message.IMessage) error
	Load([]byte, uint64) (message.IMessage, error)

	Close() error
	Clean() error
}
