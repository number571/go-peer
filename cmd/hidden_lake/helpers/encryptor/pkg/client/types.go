package client

import (
	"context"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/config"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	EncryptMessage(context.Context, asymmetric.IPubKey, []byte) (net_message.IMessage, error)
	DecryptMessage(context.Context, net_message.IMessage) (asymmetric.IPubKey, []byte, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	EncryptMessage(context.Context, asymmetric.IPubKey, []byte) (net_message.IMessage, error)
	DecryptMessage(context.Context, net_message.IMessage) (asymmetric.IPubKey, []byte, error)
}
