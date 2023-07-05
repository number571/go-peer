package storage

type ISettings interface {
	GetPath() string
	GetHashing() bool
	GetCipherKey() []byte
}

type IKVStorage interface {
	GetSettings() ISettings

	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Del([]byte) error
}
