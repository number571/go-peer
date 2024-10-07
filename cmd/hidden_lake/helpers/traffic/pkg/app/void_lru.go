package app

import (
	"github.com/number571/go-peer/pkg/storage/cache/lru"
)

var (
	_ lru.ILRUCache = &sVoidLRUCache{}
)

type sVoidLRUCache struct{}

func newVoidLRUCache() lru.ILRUCache { return &sVoidLRUCache{} }

func (p *sVoidLRUCache) GetIndex() uint64               { return 0 }
func (p *sVoidLRUCache) GetKey(_ uint64) ([]byte, bool) { return nil, false }
func (p *sVoidLRUCache) Get(_ []byte) ([]byte, bool)    { return nil, false }
func (p *sVoidLRUCache) Set(_, _ []byte) bool           { return true }
