package client

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IClient interface {
	GetIndex() (string, error)

	GetPointer() (uint64, error)
	GetHash(uint64) (string, error)

	GetMessage(string) (net_message.IMessage, error)
	PutMessage(net_message.IMessage) error
}

type IBuilder interface {
	PutMessage(net_message.IMessage) string
}

type IRequester interface {
	GetIndex() (string, error)

	GetPointer() (uint64, error)
	GetHash(uint64) (string, error)

	GetMessage(string) (net_message.IMessage, error)
	PutMessage(string) error
}
