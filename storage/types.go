package storage

type IKeyValueStorage interface {
	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Del([]byte) error
	Close() error
}
