package settings

import "github.com/number571/go-peer/pkg/client/message"

type IConfigSettingsHLS interface {
	IConfigSettings

	GetQueuePeriodMS() uint64
	GetKeySizeBits() uint64
}

type IConfigSettingsHLT interface {
	IConfigSettings

	GetMessagesCapacity() uint64
}

type IConfigSettingsHLM interface {
	IConfigSettings

	GetKeySizeBits() uint64
}

type IConfigSettings interface {
	message.ISettings
}
