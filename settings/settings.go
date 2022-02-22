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
	_ Settings = &SettingsT{}
)

type SettingsT struct {
	mutex   sync.Mutex
	mapping map[Key]Value
}

func NewSettings() Settings {
	return &SettingsT{
		mapping: defaultSettings(),
	}
}

func (s *SettingsT) Set(k Key, v Value) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.mapping[k] = v
}

func (s *SettingsT) Get(k Key) Value {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.mapping[k]
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
