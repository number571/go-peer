package types

type IApp interface {
	IRunner
	ICloser
}

type IRunner interface {
	Run() error
}

type ICloser interface {
	Close() error
}
