package database

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/types"
)

type IDatabase interface {
	types.ICloser

	Push(net_message.IMessage) error
	Load([]byte) (net_message.IMessage, error)
	Hash(uint64) ([]byte, error)
	Pointer() uint64

	Settings() ISettings
}

type ISettings interface {
	net_message.ISettings

	GetPath() string
	GetMessagesCapacity() uint64
}

type IDBWrapper interface {
	types.ICloser

	Get() IDatabase
	Set(IDatabase) IDBWrapper
}
