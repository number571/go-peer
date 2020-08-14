package gopeer

type SettingsType map[string]interface{}
type settingsStruct struct {
	END_BYTES string
	WAIT_TIME uint
	BUFF_SIZE uint
	PACK_SIZE uint
	MAPP_SIZE uint
	AKEY_SIZE uint
	SKEY_SIZE uint
	RAND_SIZE uint
}

var settings = defaultSettings()

func defaultSettings() settingsStruct {
	return settingsStruct{
		END_BYTES: "\000\005\007\001\001\007\005\000",
		WAIT_TIME: 5,       // seconds
		BUFF_SIZE: 4 << 10, // 4KiB
		PACK_SIZE: 2 << 20, // 2MiB
		MAPP_SIZE: 1024,    // elems
		AKEY_SIZE: 1024,    // bits
		SKEY_SIZE: 16,      // bytes
		RAND_SIZE: 16,      // bytes
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
	case "WAIT_TIME":
		return settings.WAIT_TIME
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
	default:
		return 1
	}
	return 0
}

func intSettings(name string, data interface{}) uint8 {
	result := data.(uint)
	switch name {
	case "WAIT_TIME":
		settings.WAIT_TIME = result
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
