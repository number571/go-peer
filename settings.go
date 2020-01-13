package gopeer

type SettingsType map[string]interface{}
type settingsStruct struct {
	TITLE_LASTHASH   string
	TITLE_CONNECT    string
	TITLE_DISCONNECT string
	OPTION_GET       string
	OPTION_SET       string
	NETWORK          string
	VERSION          string
	BUFFSIZE         uint16
	DIFFICULTY       uint8
	RETRY_NUMB       uint8
	RETRY_TIME       uint8
	TEMPLATE         string
	HMACKEY          string
	GENESIS          string
	NOISE            string
}

var settings = defaultSettings()

func defaultSettings() settingsStruct {
	return settingsStruct{
		TITLE_LASTHASH:   "[TITLE-LASTHASH]",
		TITLE_CONNECT:    "[TITLE-CONNECT]",
		TITLE_DISCONNECT: "[TITLE-DISCONNECT]",
		OPTION_GET:       "[OPTION-GET]", // Send
		OPTION_SET:       "[OPTION-SET]", // Receive
		NETWORK:          "NETWORK-NAME",
		VERSION:          "Version 1.0.0",
		BUFFSIZE:         512,
		DIFFICULTY:       15,
		RETRY_NUMB:       3,
		RETRY_TIME:       5, // Seconds
		TEMPLATE:         "0.0.0.0",
		HMACKEY:          "PASSWORD",
		GENESIS:          "[GENESIS-PACKAGE]",
		NOISE:            "1234567890ABCDEFGHIJKLMNOPQRSTUV",
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
		case int:
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
	case "TITLE_LASTHASH":
		return settings.TITLE_LASTHASH
	case "TITLE_CONNECT":
		return settings.TITLE_CONNECT
	case "TITLE_DISCONNECT":
		return settings.TITLE_DISCONNECT
	case "OPTION_GET":
		return settings.OPTION_GET
	case "OPTION_SET":
		return settings.OPTION_SET
	case "NETWORK":
		return settings.NETWORK
	case "VERSION":
		return settings.VERSION
	case "TEMPLATE":
		return settings.TEMPLATE
	case "HMACKEY":
		return settings.HMACKEY
	case "GENESIS":
		return settings.GENESIS
	case "NOISE":
		return settings.NOISE
	case "BUFFSIZE":
		return settings.BUFFSIZE
	case "DIFFICULTY":
		return settings.DIFFICULTY
	case "RETRY_NUMB":
		return settings.RETRY_NUMB
	case "RETRY_TIME":
		return settings.RETRY_TIME
	default:
		return nil
	}
}

func stringSettings(name string, data interface{}) uint8 {
	result := data.(string)
	switch name {
	case "TITLE_LASTHASH":
		settings.TITLE_LASTHASH = result
	case "TITLE_CONNECT":
		settings.TITLE_CONNECT = result
	case "TITLE_DISCONNECT":
		settings.TITLE_DISCONNECT = result
	case "OPTION_GET":
		settings.OPTION_GET = result
	case "OPTION_SET":
		settings.OPTION_SET = result
	case "NETWORK":
		settings.NETWORK = result
	case "VERSION":
		settings.VERSION = result
	case "TEMPLATE":
		settings.TEMPLATE = result
	case "HMACKEY":
		settings.HMACKEY = result
	case "GENESIS":
		settings.GENESIS = result
	case "NOISE":
		settings.NOISE = result
	default:
		return 1
	}
	return 0
}

func intSettings(name string, data interface{}) uint8 {
	result := data.(int)
	switch name {
	case "BUFFSIZE":
		settings.BUFFSIZE = uint16(result)
	case "DIFFICULTY":
		settings.DIFFICULTY = uint8(result)
	case "RETRY_NUMB":
		settings.RETRY_NUMB = uint8(result)
	case "RETRY_TIME":
		settings.RETRY_TIME = uint8(result)
	default:
		return 1
	}
	return 0
}
