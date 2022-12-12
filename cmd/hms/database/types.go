package database

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/types"
)

type IKeyValueDB interface {
	Size([]byte) (uint64, error)
	Push([]byte, message.IMessage) error
	Load([]byte, uint64) (message.IMessage, error)

	Clean() error
	types.ICloser
}
