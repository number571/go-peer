package settings

import "sync"

const (
	MaskRout uint64 = iota + 1
	MaskPing
	TimeWait
	TimePsdo
	TimeChck
	SizePsdo
	SizeRtry
	SizeWork
	SizeConn
	SizePack
	SizeMapp
	SizeSkey
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
		MaskRout: 0xFFFFFFFFFFFFFFFF, // Include/Response package
		MaskPing: 0xEEEEEEEEEEEEEEEE, // Ping package
		TimeWait: 20,                 // seconds
		TimePsdo: 5000,               // milliseconds
		TimeChck: 60,                 // seconds
		SizePsdo: 10 << 10,           // 10*(2^10)B = 10KiB
		SizeRtry: 3,                  // quantity
		SizeWork: 20,                 // bits
		SizeConn: 10,                 // quantity
		SizePack: 8 << 20,            // 8*(2^20)B = 8MiB
		SizeMapp: 2 << 10,            // 2*(2^10)H = 88KiB
		SizeSkey: 1 << 5,             // 2^5B = 32B
	}
}
