package cache

type ILRUCache interface {
	ICache

	GetIndex() uint64
	GetKey(i uint64) ([]byte, bool)
}

type ICache interface {
	ICacheSetter
	ICacheGetter
}

type ICacheSetter interface {
	Set([]byte, interface{}) bool
}

type ICacheGetter interface {
	Get([]byte) (interface{}, bool)
}
