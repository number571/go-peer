package storage

type IKeyValueStorage interface {
	iMapping

	// some storages are not implements this methods
	// they are return default values Iter=nil, Close=nil
	Iter([]byte) iIterator
	Close() error
}

type iMapping interface {
	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Del([]byte) error
}

type iIterator interface {
	Key() []byte
	Value() []byte

	Next() bool
	Close()
}
