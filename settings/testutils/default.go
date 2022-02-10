package testutils

import (
	gp "github.com/number571/go-peer/settings"
)

func NewSettings() gp.Settings {
	settings := gp.NewSettings()
	mapping := defaultSettings()
	for k, v := range mapping {
		settings.Set(k, v)
	}
	return settings
}

// for fast tests
func defaultSettings() map[gp.Key]gp.Value {
	return map[gp.Key]gp.Value{
		gp.MaskRout: 0xFFFFFFFFFFFFFFFF, // Include/Response package
		gp.TimeWait: 50,                 // seconds
		gp.TimePsdo: 1000,               // milliseconds
		gp.SizeRtry: 1,                  // quantity
		gp.SizeWork: 10,                 // bits
		gp.SizeConn: 10,                 // quantity
		gp.SizePack: 1 << 20,            // 1*(2^20)B = 1MiB
		gp.SizeMapp: 1 << 10,            // 1*(2^10)H = 44KiB
		gp.SizeSkey: 1 << 4,             // 2^4B = 16B
	}
}
