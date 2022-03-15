package settings

import "sync"

const (
	MaskRout Key = iota + 1
	TimeWait
	TimePsdo
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
	fMapping map[Key]Value
}

func NewSettings() ISettings {
	return &sSettings{
		fMapping: defaultSettings(),
	}
}

func (s *sSettings) Set(k Key, v Value) ISettings {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	s.fMapping[k] = v
	return s
}

func (s *sSettings) Get(k Key) Value {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	v, ok := s.fMapping[k]
	if !ok {
		panic("settings: value undefined")
	}

	return v
}

func defaultSettings() map[Key]Value {
	// H - hash = len(base64(sha256(data))) = 44B
	// B - byte
	// b - bit
	return map[Key]Value{
		MaskRout: 0xFFFFFFFFFFFFFFFF, // Include/Response package
		TimeWait: 20,                 // seconds
		TimePsdo: 5000,               // milliseconds
		SizePsdo: 10 << 10,           // 10*(2^10)B = 10KiB
		SizeRtry: 3,                  // quantity
		SizeWork: 20,                 // bits
		SizeConn: 10,                 // quantity
		SizePack: 8 << 20,            // 8*(2^20)B = 8MiB
		SizeMapp: 2 << 10,            // 2*(2^10)H = 88KiB
		SizeSkey: 1 << 5,             // 2^5B = 32B
	}
}
