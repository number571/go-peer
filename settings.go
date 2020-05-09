package gopeer

type SettingsType map[string]interface{}
type settingsStruct struct {
	TITLE_CONNECT      string
	TITLE_DISCONNECT   string
	TITLE_FILETRANSFER string
	OPTION_GET         string
	OPTION_SET         string
	IS_CLIENT          string
	END_BYTES          string
	TEMPLATE           string
	HMAC_KEY           string
	NETWORK            string
	VERSION            string
	max_id             uint64
	BITS_SIZE          uint64
	PACK_SIZE          uint64
	BUFF_SIZE          uint32
	RAND_SIZE          uint16
	KEY_SIZE           uint16
	REMEMBER           uint16
	DIFFICULTY         uint8
	RETRY_QUAN         uint8
	WAITING_TIME       uint8
	SESSION_SIZE       uint8
	REDIRECT_QUAN      uint8
}

var settings = defaultSettings()

// b = bit
// B = byte
func defaultSettings() settingsStruct {
	return settingsStruct{
		TITLE_CONNECT:      "[TITLE-CONNECT]",
		TITLE_DISCONNECT:   "[TITLE-DISCONNECT]",
		TITLE_FILETRANSFER: "[TITLE-FILETRANSFER]",
		OPTION_GET:         "[OPTION-GET]", // Send
		OPTION_SET:         "[OPTION-SET]", // Receive
		IS_CLIENT:          "[IS-CLIENT]",
		END_BYTES:          "\000\000\000\005\007\001\000\000\000",
		TEMPLATE:           "0.0.0.0",
		HMAC_KEY:           "PASSWORD",
		NETWORK:            "GOPEER-FRAMEWORK",
		VERSION:            "Version 1.0.0",
		max_id:             (1 << 48) / (8 << 20), // BITS_SIZE / PACK_SIZE
		BITS_SIZE:          1 << 48,               // 2^48b
		PACK_SIZE:          8 << 20,               // 8MiB
		BUFF_SIZE:          1 << 20,               // 1MiB
		RAND_SIZE:          1 << 4,                // 16B
		KEY_SIZE:           2 << 10,               // 2048b
		REMEMBER:           256,                   // quantity hash packages
		DIFFICULTY:         15,                    // bits of 256b
		RETRY_QUAN:         2,                     // quantity retry send one package
		WAITING_TIME:       5,                     // in seconds
		SESSION_SIZE:       32,                    // bytes for AES128/192/256
		REDIRECT_QUAN:      3,                     // quantity hidden nodes that can send a package
	}
}

// 0 = success
// 1 = parameter undefined
// 2 = type undefined
func Set(settings SettingsType) []uint8 {
	var (
		list = make([]uint8, len(settings))
		i    = 0
	)

	for name, data := range settings {
		switch data.(type) {
		case string:
			list[i] = stringSettings(name, data)
		case uint64:
			list[i] = intSettings(name, data)
		default:
			list[i] = 2
		}
		i++
	}

	return list
}

func Get(key string) interface{} {
	switch key {
	case "TITLE_CONNECT":
		return settings.TITLE_CONNECT
	case "TITLE_DISCONNECT":
		return settings.TITLE_DISCONNECT
	case "TITLE_FILETRANSFER":
		return settings.TITLE_FILETRANSFER
	case "OPTION_GET":
		return settings.OPTION_GET
	case "OPTION_SET":
		return settings.OPTION_SET
	case "IS_CLIENT":
		return settings.IS_CLIENT
	case "END_BYTES":
		return settings.END_BYTES
	case "NETWORK":
		return settings.NETWORK
	case "VERSION":
		return settings.VERSION
	case "TEMPLATE":
		return settings.TEMPLATE
	case "HMAC_KEY":
		return settings.HMAC_KEY
	case "BITS_SIZE":
		return settings.BITS_SIZE
	case "PACK_SIZE":
		return settings.PACK_SIZE
	case "BUFF_SIZE":
		return settings.BUFF_SIZE
	case "RAND_SIZE":
		return settings.RAND_SIZE
	case "KEY_SIZE":
		return settings.KEY_SIZE
	case "REMEMBER":
		return settings.REMEMBER
	case "DIFFICULTY":
		return settings.DIFFICULTY
	case "RETRY_QUAN":
		return settings.RETRY_QUAN
	case "WAITING_TIME":
		return settings.WAITING_TIME
	case "SESSION_SIZE":
		return settings.SESSION_SIZE
	case "REDIRECT_QUAN":
		return settings.REDIRECT_QUAN
	default:
		return nil
	}
}

func stringSettings(name string, data interface{}) uint8 {
	result := data.(string)
	switch name {
	case "TITLE_CONNECT":
		settings.TITLE_CONNECT = result
	case "TITLE_DISCONNECT":
		settings.TITLE_DISCONNECT = result
	case "TITLE_FILETRANSFER":
		settings.TITLE_FILETRANSFER = result
	case "OPTION_GET":
		settings.OPTION_GET = result
	case "OPTION_SET":
		settings.OPTION_SET = result
	case "IS_CLIENT":
		settings.IS_CLIENT = result
	case "END_BYTES":
		settings.END_BYTES = result
	case "NETWORK":
		settings.NETWORK = result
	case "VERSION":
		settings.VERSION = result
	case "TEMPLATE":
		settings.TEMPLATE = result
	case "HMAC_KEY":
		settings.HMAC_KEY = result
	default:
		return 1
	}
	return 0
}

func intSettings(name string, data interface{}) uint8 {
	result := data.(uint64)
	switch name {
	case "BITS_SIZE":
		settings.BITS_SIZE = uint64(result)
		settings.max_id = settings.BITS_SIZE / settings.PACK_SIZE
	case "PACK_SIZE":
		settings.PACK_SIZE = uint64(result)
		settings.max_id = settings.BITS_SIZE / settings.PACK_SIZE
	case "BUFF_SIZE":
		settings.BUFF_SIZE = uint32(result)
	case "RAND_SIZE":
		settings.RAND_SIZE = uint16(result)
	case "KEY_SIZE":
		settings.KEY_SIZE = uint16(result)
	case "REMEMBER":
		settings.REMEMBER = uint16(result)
	case "DIFFICULTY":
		settings.DIFFICULTY = uint8(result)
	case "RETRY_QUAN":
		settings.RETRY_QUAN = uint8(result)
	case "WAITING_TIME":
		settings.WAITING_TIME = uint8(result)
	case "SESSION_SIZE":
		size := uint8(result)
		switch size {
		case 16, 24, 32:
			settings.SESSION_SIZE = size
		default:
			return 1
		}
	case "REDIRECT_QUAN":
		settings.REDIRECT_QUAN = uint8(result)
	default:
		return 1
	}
	return 0
}
