package config

import "github.com/number571/go-peer/pkg/client/message"

type IConfigSettings interface {
	message.ISettings

	GetKeySizeBits() uint64
	GetQueuePeriodMS() uint64
	GetLimitVoidSizeBytes() uint64
}
