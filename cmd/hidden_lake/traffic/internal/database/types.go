package database

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/types"
)

type IKVDatabase interface {
	Push(message.IMessage) error
	Load(string) (message.IMessage, error)
	Hashes() ([]string, error)

	Settings() ISettings
	types.ICloser
}

type ISettings interface {
	GetPath() string
	GetCapacity() uint64
	GetMessageSize() uint64
	GetWorkSize() uint64
}

type IWrapperDB interface {
	types.ICloser

	Get() IKVDatabase
	Set(IKVDatabase) IWrapperDB
}
