package client

import (
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/config"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IClient interface {
	GetIndex() (string, error)
	GetSettings() (config.IConfigSettings, error)

	GetPubKey() (asymmetric.IPubKey, error)

	EncryptMessage(asymmetric.IPubKey, []byte) (net_message.IMessage, error)
	DecryptMessage(net_message.IMessage) (asymmetric.IPubKey, []byte, error)
}

type IRequester interface {
	GetIndex() (string, error)
	GetSettings() (config.IConfigSettings, error)

	GetPubKey() (asymmetric.IPubKey, error)

	EncryptMessage(asymmetric.IPubKey, []byte) (net_message.IMessage, error)
	DecryptMessage(net_message.IMessage) (asymmetric.IPubKey, []byte, error)
}