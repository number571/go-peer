package types

import (
	"context"
	"io"
)

type ICloser io.Closer

type IRunner interface {
	Run(context.Context) error
}

type IConverter interface {
	ToString() string
	ToBytes() []byte
}
