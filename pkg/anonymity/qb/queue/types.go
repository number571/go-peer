package queue

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/types"
)

type IQBProblemProcessor interface {
	types.IRunner

	GetSettings() ISettings
	GetScheme() layer2.IScheme

	EnqueueMessage(layer2.IParticipantKey, []byte) error
	DequeueMessage(context.Context) layer1.IMessage
}

type ISettings interface {
	GetMessageConstructSettings() layer1.IConstructSettings
	GetNetworkMask() uint32
	GetConsumersCap() uint64
	GetQueuePeriod() time.Duration
	GetQueuePoolCap() [2]uint64
}
