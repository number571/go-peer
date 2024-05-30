package client

import (
	"context"

	"github.com/number571/go-peer/cmd/hidden_lake/template/pkg/config"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	// ...
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	// ...
}
