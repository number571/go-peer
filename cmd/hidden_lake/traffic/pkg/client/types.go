package client

import (
	"github.com/number571/go-peer/pkg/client/message"
)

type IClient interface {
	GetIndex() (string, error)

	GetHashes() ([][]byte, error)

	GetMessage(string) (message.IMessage, error)
	PutMessage(message.IMessage) error
}

type IBuilder interface {
	PutMessage(message.IMessage) string
}

type IRequester interface {
	GetIndex() (string, error)

	GetHashes() ([][]byte, error)

	GetMessage(string) (message.IMessage, error)
	PutMessage(string) error
}
