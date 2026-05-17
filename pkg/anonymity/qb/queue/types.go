package queue

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	hybrid "github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/types"
)

type IQBProblemProcessor interface {
	types.IRunner

	GetSettings() ISettings
	GetScheme() hybrid.IScheme

	EnqueueMessage(hybrid.IParticipantKey, []byte) error
	DequeueMessage(context.Context) layer1.IMessage
}

type ISettings interface {
	GetMessageConstructSettings() layer1.IConstructSettings
	GetNetworkMask() uint32
	GetConsumersCap() uint64
	GetQueuePeriod() time.Duration
	GetQueuePoolCap() [2]uint64
}
