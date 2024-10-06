package storage

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/cache/lru"
	"github.com/number571/go-peer/pkg/storage/database"
)

type IMessageStorage interface {
	GetSettings() net_message.ISettings
	GetKVDatabase() database.IKVDatabase
	GetLRUCache() lru.ILRUCache

	Pointer() uint64
	Push(net_message.IMessage) error
	Load([]byte) (net_message.IMessage, error)
	Hash(uint64) ([]byte, error)
}
