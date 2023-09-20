package storage

type IKVStorage interface {
	GetSettings() ISettings

	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Del([]byte) error
}

type ISettings interface {
	GetPath() string
	GetWorkSize() uint64
	GetPassword() string
}
