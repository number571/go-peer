package settings

import "sync"

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

func (s *SettingsT) Set(k Key, v Value) Settings {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.mapping[k] = v
	return s
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

// func defaultSettings() map[Key]Value {
// 	// H - hash = len(base64(sha256(data))) = 44B
// 	// B - byte
// 	// b - bit
// 	return map[Key]Value{
// 		MaskRout: 0xFFFFFFFFFFFFFFFF, // Include/Response package
// 		TimeWait: 20,                 // seconds
// 		TimePsdo: 5000,               // milliseconds
// 		SizeRtry: 3,                  // quantity
// 		SizeWork: 20,                 // bits
// 		SizeConn: 10,                 // quantity
// 		SizePack: 8 << 20,            // 8*(2^20)B = 8MiB
// 		SizeMapp: 2 << 10,            // 2*(2^10)H = 88KiB
// 		SizeAkey: 2 << 10,            // 2*(2^10)b = 256B
// 		SizeSkey: 1 << 5,             // 2^5B = 32B
// 	}
// }

// for fast tests
func defaultSettings() map[Key]Value {
	return map[Key]Value{
		MaskRout: 0xFFFFFFFFFFFFFFFF, // Include/Response package
		TimeWait: 10,                 // seconds
		TimePsdo: 1000,               // milliseconds
		SizeRtry: 1,                  // quantity
		SizeWork: 10,                 // bits
		SizeConn: 10,                 // quantity
		SizePack: 2 << 20,            // 2*(2^20)B = 2MiB
		SizeMapp: 1 << 10,            // 1*(2^10)H = 44KiB
		SizeAkey: 1 << 10,            // 1*(2^10)b = 128B
		SizeSkey: 1 << 4,             // 2^4B = 16B
	}
}
