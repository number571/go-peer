package cache

type ILRUCache interface {
	ICache

	GetIndex() uint64
	GetKey(i uint64) (string, bool)
}

type ICache interface {
	ICacheSetter
	ICacheGetter
}

type ICacheSetter interface {
	Set(string, interface{}) bool
}

type ICacheGetter interface {
	Get(string) (interface{}, bool)
}
