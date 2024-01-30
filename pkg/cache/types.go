package cache

type ICacheSetter interface {
	Set([]byte, []byte) bool
}

type ICacheGetter interface {
	Get([]byte) ([]byte, bool)
}
