package database

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/types"
)

type IKeyValueDB interface {
	Push(message.IMessage) error
	Load(string) (message.IMessage, error)
	Hashes() ([]string, error)

	Settings() ISettings
	types.ICloser
}

type ISettings interface {
	GetPath() string
	GetLimitMessages() uint64
	GetMessageSize() uint64
	GetWorkSize() uint64
}
