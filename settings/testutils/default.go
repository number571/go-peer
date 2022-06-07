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
		settings.CMaskRout: 0xFFFFFFFFFFFFFFFF, // Include/Response package
		settings.CMaskPing: 0xEEEEEEEEEEEEEEEE, // Ping package
		settings.CTimeWait: 20,                 // seconds
		settings.CTimePreq: 1,                  // seconds
		settings.CTimePrsp: 1,                  // seconds
		settings.CTimePing: 1,                  // seconds
		settings.CSizePsdo: 2 << 10,            // 2*(2^10)B = 2KiB
		settings.CSizeRtry: 0,                  // quantity
		settings.CSizeWork: 10,                 // bits
		settings.CSizeConn: 10,                 // quantity
		settings.CSizePack: 1 << 20,            // 1*(2^20)B = 1MiB
		settings.CSizeMapp: 1 << 10,            // 1*(2^10)H = 44KiB
		settings.CSizeSkey: 1 << 4,             // 2^4B = 16B
		settings.CSizeBmsg: 10,                 // quantity
	}
}
