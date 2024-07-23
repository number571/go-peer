package lru

import "github.com/number571/go-peer/pkg/storage/cache"

type ILRUCache interface {
	cache.ICache

	GetIndex() uint64
	GetKey(i uint64) ([]byte, bool)
}
