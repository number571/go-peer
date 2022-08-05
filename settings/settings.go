package settings

import (
	"sync"
)

var (
	_ ISettings = &sSettings{}
)

type sSettings struct {
	fMutex   sync.Mutex
	fMapping map[uint64]uint64
}

func NewSettings() ISettings {
	return &sSettings{
		fMapping: defaultSettings(),
	}
}

func (s *sSettings) Set(k uint64, v uint64) ISettings {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	s.fMapping[k] = v
	return s
}

func (s *sSettings) Get(k uint64) uint64 {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	v, ok := s.fMapping[k]
	if !ok {
		panic("settings: value undefined")
	}

	return v
}

// Need to change the default settings!
func defaultSettings() map[uint64]uint64 {
	// H - hash = len(base64(sha256(data))) = 44B
	// B - byte
	// b - bit
	return map[uint64]uint64{
		CMaskNetw: 0xFFFFFFFFFFFFFFFF, // Network route
		CMaskRout: 0xEEEEEEEEEEEEEEEE, // Include/Response package
		CMaskPing: 0x00000000DDDDDDDD, // Must be 32bit; Ping package
		CMaskPsdo: 0x0000000000000000, // Used for pseudo packages
		CMaskPasw: 0b000,              // 0b111 = (alpha, numeric, special)
		CTimeWait: 20,                 // seconds
		CTimePreq: 1,                  // seconds (used for pseudo requests)
		CTimeRslp: 1,                  // seconds (used for random sleep)
		CTimePing: 1,                  // seconds (used for online checker)
		CSizeRout: 3,                  // Must be > 0; Max routes
		CSizePsdo: 2 << 10,            // 2*(2^10)B = 2KiB
		CSizeRtry: 0,                  // quantity
		CSizeWork: 10,                 // bits
		CSizeConn: 10,                 // quantity
		CSizePack: 1 << 20,            // 1*(2^20)B = 1MiB
		CSizeMapp: 1 << 10,            // quantity hashes in map
		CSizeSkey: 1 << 4,             // 2^4B = 16B
		CSizeBmsg: 10,                 // quantity
		CSizePasw: 4,                  // chars
	}
}
