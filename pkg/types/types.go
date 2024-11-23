package types

import (
	"context"
)

type IRunner interface {
	Run(context.Context) error
}

type IConverter interface {
	ToString() string
	ToBytes() []byte
}
