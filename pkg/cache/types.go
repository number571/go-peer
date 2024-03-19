package cache

type ICache interface {
	ICacheSetter
	ICacheGetter
}

type ICacheSetter interface {
	Set([]byte, []byte) bool
}

type ICacheGetter interface {
	Get([]byte) ([]byte, bool)
}
