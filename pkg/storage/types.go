package storage

type ISettings interface {
	GetPath() string
	GetWorkSize() uint64
	GetCipherKey() []byte
}

type IKeyValueStorage interface {
	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Del([]byte) error
}
