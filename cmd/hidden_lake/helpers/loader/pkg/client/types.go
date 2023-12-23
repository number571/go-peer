package client

import "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/config"

type IClient interface {
	GetIndex() (string, error)
	GetSettings() (config.IConfigSettings, error)

	RunTransfer() error
	StopTransfer() error
}

type IRequester interface {
	GetIndex() (string, error)
	GetSettings() (config.IConfigSettings, error)

	RunTransfer() error
	StopTransfer() error
}
