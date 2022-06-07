package settings

import "sync"

const (
	CMaskRout uint64 = iota + 1
	CMaskPing
	CTimeWait
	CTimePreq
	CTimePrsp
	CTimePing
	CSizePsdo
	CSizeRtry
	CSizeWork
	CSizeConn
	CSizePack
	CSizeMapp
	CSizeSkey
	CSizeBmsg
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

func defaultSettings() map[uint64]uint64 {
	// H - hash = len(base64(sha256(data))) = 44B
	// B - byte
	// b - bit
	return map[uint64]uint64{
		CMaskRout: 0xFFFFFFFFFFFFFFFF, // Include/Response package
		CMaskPing: 0xEEEEEEEEEEEEEEEE, // Ping package
		CTimeWait: 30,                 // seconds
		CTimePreq: 10,                 // seconds
		CTimePrsp: 5,                  // seconds
		CTimePing: 60,                 // seconds
		CSizePsdo: 10 << 10,           // 10*(2^10)B = 10KiB
		CSizeRtry: 2,                  // quantity
		CSizeWork: 20,                 // bits
		CSizeConn: 20,                 // quantity
		CSizePack: 8 << 20,            // 8*(2^20)B = 8MiB
		CSizeMapp: 2 << 10,            // 2*(2^10)H = 88KiB
		CSizeSkey: 1 << 5,             // 2^5B = 32B
		CSizeBmsg: 20,                 // quantity of messages
	}
}
