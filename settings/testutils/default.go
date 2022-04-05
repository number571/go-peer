package testutils

import (
	"github.com/number571/go-peer/settings"
)

func NewSettings() settings.ISettings {
	sett := settings.NewSettings()
	mapping := defaultSettings()
	for k, v := range mapping {
		sett.Set(k, v)
	}
	return sett
}

// for fast tests
func defaultSettings() map[uint64]uint64 {
	return map[uint64]uint64{
		settings.MaskRout: 0xFFFFFFFFFFFFFFFF, // Include/Response package
		settings.TimeWait: 50,                 // seconds
		settings.TimePsdo: 1000,               // milliseconds
		settings.SizeRtry: 1,                  // quantity
		settings.SizeWork: 10,                 // bits
		settings.SizeConn: 10,                 // quantity
		settings.SizePack: 1 << 20,            // 1*(2^20)B = 1MiB
		settings.SizeMapp: 1 << 10,            // 1*(2^10)H = 44KiB
		settings.SizeSkey: 1 << 4,             // 2^4B = 16B
	}
}
