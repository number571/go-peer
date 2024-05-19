package queue

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/types"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IMessageQueue interface {
	types.IRunner

	SetVSettings(IVSettings)
	GetVSettings() IVSettings

	GetSettings() ISettings
	GetClient() client.IClient

	EnqueueMessage(message.IMessage) error
	DequeueMessage(context.Context) net_message.IMessage
}

type ISettings interface {
	GetNetworkMask() uint64
	GetWorkSizeBits() uint64
	GetMainCapacity() uint64
	GetVoidCapacity() uint64
	GetParallel() uint64
	GetDuration() time.Duration
	GetRandDuration() time.Duration
	GetLimitVoidSizeBytes() uint64
}

type IVSettings interface {
	GetNetworkKey() string
}
