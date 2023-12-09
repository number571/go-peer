package database

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/types"
)

type IDatabase interface {
	types.ICloser

	Push(net_message.IMessage) error
	Load([]byte) (net_message.IMessage, error)
	Hashes() ([][]byte, error)

	Settings() ISettings
}

type ISettings interface {
	net_message.ISettings

	GetPath() string
	GetHashesWindow() uint64
	GetMessagesCapacity() uint64
}

type IWrapperDB interface {
	types.ICloser

	Get() IDatabase
	Set(IDatabase) IWrapperDB
}
