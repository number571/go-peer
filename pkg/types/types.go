package types

type ICommand interface {
	Run() error
	Stop() error
}

type ICloser interface {
	Close() error
}

type IConverter interface {
	ToString() string
	ToBytes() []byte
}
