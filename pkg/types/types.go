package types

import (
	"context"
	"io"
)

type IRunner interface {
	Run(context.Context) error
}

type ICloser io.Closer

type IConverter interface {
	ToString() string
	ToBytes() []byte
}
