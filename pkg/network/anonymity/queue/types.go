package queue

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IMessageQueueProcessor interface {
	types.IRunner

	SetVSettings(IVSettings)
	GetVSettings() IVSettings

	GetSettings() ISettings
	GetClient() client.IClient

	EnqueueMessage(asymmetric.IPubKey, []byte) error
	DequeueMessage(context.Context) net_message.IMessage
}

type ISettings interface {
	GetNetworkMask() uint32
	GetWorkSizeBits() uint64
	GetParallel() uint64
	GetMainPoolCapacity() uint64
	GetRandPoolCapacity() uint64
	GetQueuePeriod() time.Duration
	GetRandQueuePeriod() time.Duration
	GetRandMessageSizeBytes() uint64
}

type IVSettings interface {
	GetNetworkKey() string
}
