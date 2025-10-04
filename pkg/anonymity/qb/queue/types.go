package queue

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/types"
)

type IQBProblemProcessor interface {
	types.IRunner

	GetSettings() ISettings
	GetClient() client.IClient

	EnqueueMessage(asymmetric.IPubKey, []byte) error
	DequeueMessage(context.Context) layer1.IMessage
}

type ISettings interface {
	GetMessageConstructSettings() layer1.IConstructSettings
	GetNetworkMask() uint32
	GetConsumersCap() uint64
	GetQueuePeriod() time.Duration
	GetQueuePoolCap() [2]uint64
}
