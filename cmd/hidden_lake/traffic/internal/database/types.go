package database

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/types"
)

type IKVDatabase interface {
	types.ICloser

	Push(message.IMessage) error
	Load(string) (message.IMessage, error)
	Hashes() ([]string, error)

	Settings() ISettings
}

type ISettings interface {
	GetPath() string
	GetMessagesCapacity() uint64
	GetMessageSizeBytes() uint64
	GetWorkSizeBits() uint64
}

type IWrapperDB interface {
	types.ICloser

	Get() IKVDatabase
	Set(IKVDatabase) IWrapperDB
}
