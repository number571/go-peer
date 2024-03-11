package client

import (
	"context"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/config"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPointer(context.Context) (uint64, error)
	GetHash(context.Context, uint64) (string, error)

	GetMessage(context.Context, string) (net_message.IMessage, error)
	PutMessage(context.Context, net_message.IMessage) error
}

type IBuilder interface {
	PutMessage(net_message.IMessage) string
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPointer(context.Context) (uint64, error)
	GetHash(context.Context, uint64) (string, error)

	GetMessage(context.Context, string) (net_message.IMessage, error)
	PutMessage(context.Context, string) error
}
