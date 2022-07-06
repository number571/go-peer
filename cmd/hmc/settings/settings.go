package settings

import "github.com/number571/go-peer/settings"

func NewSettings() settings.ISettings {
	// another parameters are not used
	return settings.NewSettings().
		Set(settings.CSizeSkey, 1<<5).              // bytes
		Set(settings.CSizeWork, 25).                // bits
		Set(settings.CSizePasw, 8).                 // chars
		Set(settings.CMaskPasw, fullPasswordMode()) // passwords rule
}

func fullPasswordMode() uint64 {
	return settings.CPaswAplh | settings.CPaswNumr | settings.CPaswSpec
}
