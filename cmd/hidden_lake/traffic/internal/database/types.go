package database

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
)

type IKVDatabase interface {
	types.ICloser

	Push(message.IMessage) error
	Load(string) (message.IMessage, error)
	Hashes() ([]string, error)

	Settings() ISettings
	GetOriginal() database.IKVDatabase
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
