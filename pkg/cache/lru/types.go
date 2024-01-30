package lru

import "github.com/number571/go-peer/pkg/cache"

type ILRUCache interface {
	GetSettings() ISettings

	GetIndex() uint64
	GetKey(i uint64) ([]byte, bool)

	cache.ICacheSetter
	cache.ICacheGetter
}

type ISettings interface {
	GetCapacity() uint64
}
