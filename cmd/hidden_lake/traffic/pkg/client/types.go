package client

import (
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/client/message"
)

type IClient interface {
	GetIndex() (string, error)
	GetHashes() ([]string, error)
	GetMessage(string) (message.IMessage, error)
	AddMessage(message.IMessage) error
	DoBroadcast() error
}

type IBuilder interface {
	GetMessage(string) *pkg_settings.SLoadRequest
	AddMessage(message.IMessage) *pkg_settings.SPushRequest
}

type IRequester interface {
	GetIndex() (string, error)
	GetHashes() ([]string, error)
	GetMessage(*pkg_settings.SLoadRequest) (message.IMessage, error)
	AddMessage(*pkg_settings.SPushRequest) error
	DoBroadcast() error
}
