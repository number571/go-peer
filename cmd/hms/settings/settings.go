package settings

import "github.com/number571/go-peer/settings"

func NewSettings() settings.ISettings {
	// another parameters are not used
	return settings.NewSettings().
		Set(settings.CSizePack, 8<<20). // 8MiB
		Set(settings.CSizeWork, 25)     // 25bits
}
