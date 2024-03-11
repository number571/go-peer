package client

import (
	"context"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/config"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	RunTransfer(context.Context) error
	StopTransfer(context.Context) error
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	RunTransfer(context.Context) error
	StopTransfer(context.Context) error
}
