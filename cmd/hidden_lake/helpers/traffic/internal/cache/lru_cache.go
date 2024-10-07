package cache

import (
	"github.com/number571/go-peer/pkg/storage/cache"
)

func NewLRUCache(pCapacity uint64) cache.ILRUCache {
	if pCapacity == 0 {
		return newVoidLRUCache()
	}
	return cache.NewLRUCache(pCapacity)
}
