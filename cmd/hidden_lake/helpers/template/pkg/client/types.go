package client

import "github.com/number571/go-peer/cmd/hidden_lake/helpers/template/pkg/config"

type IClient interface {
	GetIndex() (string, error)
	GetSettings() (config.IConfigSettings, error)

	// TODO: need implementation
}

type IRequester interface {
	GetIndex() (string, error)
	GetSettings() (config.IConfigSettings, error)

	// TODO: need implementation
}
