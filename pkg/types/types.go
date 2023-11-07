package types

type ICommand interface {
	IRunner
	IStopper
}

type IRunner interface {
	Run() error
}

type IStopper interface {
	Stop() error
}

type ICloser interface {
	Close() error
}

type IConverter interface {
	ToString() string
	ToBytes() []byte
}
