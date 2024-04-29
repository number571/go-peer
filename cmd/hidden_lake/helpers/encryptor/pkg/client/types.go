package client

import (
	"context"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/config"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	EncryptMessage(context.Context, asymmetric.IPubKey, payload.IPayload) (net_message.IMessage, error)
	DecryptMessage(context.Context, net_message.IMessage) (asymmetric.IPubKey, payload.IPayload, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	EncryptMessage(context.Context, asymmetric.IPubKey, payload.IPayload) (net_message.IMessage, error)
	DecryptMessage(context.Context, net_message.IMessage) (asymmetric.IPubKey, payload.IPayload, error)
}
