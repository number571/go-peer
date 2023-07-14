package settings

import "github.com/number571/go-peer/pkg/client/message"

type IConfigSettingsHLS interface {
	IConfigSettings

	GetQueuePeriod() uint64
	GetKeySize() uint64
}

type IConfigSettingsHLT interface {
	IConfigSettings

	GetCapMessages() uint64
}

type IConfigSettingsHLM interface {
	IConfigSettings

	GetKeySize() uint64
}

type IConfigSettings interface {
	message.ISettings
}
