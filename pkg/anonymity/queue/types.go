package queue

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IQBProblemProcessor interface {
	types.IRunner

	GetSettings() ISettings
	GetClient() client.IClient

	EnqueueMessage(asymmetric.IPubKey, []byte) error
	DequeueMessage(context.Context) net_message.IMessage
}

type ISettings interface {
	GetMessageConstructSettings() net_message.IConstructSettings
	GetNetworkMask() uint32
	GetConsumersCap() uint64
	GetQueuePeriod() time.Duration
	GetQueuePoolCap() [2]uint64
}
