package queue

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/types"
)

type IMessageQueue interface {
	types.IRunner
	SetNetworkSettings(uint64, string)

	GetSettings() ISettings
	GetClient() client.IClient

	EnqueueMessage(message.IMessage) error
	DequeueMessage(context.Context) net_message.IMessage
}

type ISettings interface {
	net_message.ISettings
	GetNetworkMask() uint64

	GetMainCapacity() uint64
	GetVoidCapacity() uint64
	GetParallel() uint64
	GetDuration() time.Duration
	GetRandDuration() time.Duration
	GetLimitVoidSizeBytes() uint64
}
