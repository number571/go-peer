package client

import (
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/client/message"
)

type IClient interface {
	Hashes() ([]string, error)
	Load(string) (message.IMessage, error)
	Push(message.IMessage) error
}

type IBuilder interface {
	Load(string) *pkg_settings.SLoadRequest
	Push(message.IMessage) *pkg_settings.SPushRequest
}

type IRequester interface {
	Hashes() ([]string, error)
	Load(*pkg_settings.SLoadRequest) (message.IMessage, error)
	Push(*pkg_settings.SPushRequest) error
}
