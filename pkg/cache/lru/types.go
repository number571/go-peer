package lru

import "github.com/number571/go-peer/pkg/cache"

type ILRUCache interface {
	cache.ICache
	GetSettings() ISettings

	GetIndex() uint64
	GetKey(i uint64) ([]byte, bool)
}

type ISettings interface {
	GetCapacity() uint64
}
