package gopeer

type SettingsType map[string]interface{}
type settingsStruct struct {
	END_BYTES string
	ROUTE_MSG string
	RETRY_NUM uint
	WAIT_TIME uint
	POWS_DIFF uint
	CONN_SIZE uint
	BUFF_SIZE uint
	PACK_SIZE uint
	MAPP_SIZE uint
	AKEY_SIZE uint
	SKEY_SIZE uint
	RAND_SIZE uint
}

var settings = defaultSettings()

// H - hash = len(base64(sha256(data))) = 44B
// B - byte
// b - bit
func defaultSettings() settingsStruct {
	return settingsStruct{
		END_BYTES: "\000\005\007\001\001\007\005\000",
		ROUTE_MSG: "\000\001\002\003\004\005\006\007",
		RETRY_NUM: 3,       // quantity
		WAIT_TIME: 20,      // seconds
		POWS_DIFF: 20,      // bits
		CONN_SIZE: 10,      // quantity
		BUFF_SIZE: 2 << 20, // 2*(2^20)B = 2MiB
		PACK_SIZE: 4 << 20, // 4*(2^20)B = 4MiB
		MAPP_SIZE: 2 << 10, // 2*(2^10)H = 88KiB
		AKEY_SIZE: 2 << 10, // 2*(2^10)b = 256B
		SKEY_SIZE: 1 << 5,  // 2^5B = 32B
		RAND_SIZE: 1 << 4,  // 2^4B = 16B
	}
}

func Set(settings SettingsType) []uint8 {
	var (
		list = make([]uint8, len(settings))
		i    = 0
	)
	for name, data := range settings {
		switch data.(type) {
		case string:
			list[i] = stringSettings(name, data)
		case uint:
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
	case "END_BYTES":
		return settings.END_BYTES
	case "ROUTE_MSG":
		return settings.ROUTE_MSG
	case "RETRY_NUM":
		return settings.RETRY_NUM
	case "WAIT_TIME":
		return settings.WAIT_TIME
	case "POWS_DIFF":
		return settings.POWS_DIFF
	case "CONN_SIZE":
		return settings.CONN_SIZE
	case "BUFF_SIZE":
		return settings.BUFF_SIZE
	case "PACK_SIZE":
		return settings.PACK_SIZE
	case "MAPP_SIZE":
		return settings.MAPP_SIZE
	case "AKEY_SIZE":
		return settings.AKEY_SIZE
	case "SKEY_SIZE":
		return settings.SKEY_SIZE
	case "RAND_SIZE":
		return settings.RAND_SIZE
	default:
		return nil
	}
}

func stringSettings(name string, data interface{}) uint8 {
	result := data.(string)
	switch name {
	case "END_BYTES":
		settings.END_BYTES = result
	case "ROUTE_MSG":
		settings.ROUTE_MSG = result
	default:
		return 1
	}
	return 0
}

func intSettings(name string, data interface{}) uint8 {
	result := data.(uint)
	switch name {
	case "RETRY_NUM":
		settings.RETRY_NUM = result
	case "WAIT_TIME":
		settings.WAIT_TIME = result
	case "POWS_DIFF":
		settings.POWS_DIFF = result
	case "CONN_SIZE":
		settings.CONN_SIZE = result
	case "BUFF_SIZE":
		settings.BUFF_SIZE = result
	case "PACK_SIZE":
		settings.PACK_SIZE = result
	case "MAPP_SIZE":
		settings.MAPP_SIZE = result
	case "AKEY_SIZE":
		settings.AKEY_SIZE = result
	case "SKEY_SIZE":
		settings.SKEY_SIZE = result
	case "RAND_SIZE":
		settings.RAND_SIZE = result
	default:
		return 1
	}
	return 0
}
