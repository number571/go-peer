package types

import "context"

type IRunner interface {
	Run(context.Context) error
}

type ICloser interface {
	Close() error
}

type IConverter interface {
	ToString() string
	ToBytes() []byte
}
