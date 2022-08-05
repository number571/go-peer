package settings

import "github.com/number571/go-peer/settings"

func NewSettings() settings.ISettings {
	sett := settings.NewSettings()
	mapping := defaultSettings()
	for k, v := range mapping {
		sett.Set(k, v)
	}
	return sett
}

func defaultSettings() map[uint64]uint64 {
	// H - hash = len(base64(sha256(data))) = 44B
	// B - byte
	// b - bit
	return map[uint64]uint64{
		settings.CMaskNetw: 0xFFFFFFFFFFFFFFFF, // Network route
		settings.CMaskRout: 0xEEEEEEEEEEEEEEEE, // Include/Response package
		settings.CMaskPing: 0x00000000DDDDDDDD, // Must be 32bit; Ping package
		settings.CMaskPsdo: 0x0000000000000000, // Used for pseudo packages
		settings.CMaskPasw: fullPasswordMode(), // 0b111 = (alpha, numeric, special)
		settings.CTimeWait: 60,                 // seconds
		settings.CTimePreq: 10,                 // seconds
		settings.CTimeRslp: 5,                  // seconds
		settings.CTimePing: 60,                 // seconds
		settings.CSizeRout: 3,                  // max routes
		settings.CSizePsdo: 10 << 10,           // 10*(2^10)B = 10KiB
		settings.CSizeRtry: 2,                  // quantity
		settings.CSizeWork: 20,                 // bits
		settings.CSizeConn: 20,                 // quantity
		settings.CSizePack: 8 << 20,            // 8*(2^20)B = 8MiB
		settings.CSizeMapp: 2 << 10,            // 2*(2^10)H = 88KiB
		settings.CSizeSkey: 1 << 5,             // 2^5B = 32B
		settings.CSizeBmsg: 20,                 // quantity of messages
		settings.CSizePasw: 4,                  // chars
	}
}

func fullPasswordMode() uint64 {
	return settings.CPaswAplh | settings.CPaswNumr | settings.CPaswSpec
}
