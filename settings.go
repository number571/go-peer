package gopeer

import (
    "time"
)

const (
    // len(base64(sha256(X))) = 24
    LEN_BASE64_SHA256 = 24
    DATA_SIZE = 3
)

var setting = struct {
    TEMPLATE string
    CLIENT_NAME string
    SEPARATOR string
    END_BYTES string
    NETWORK_NAME string
    TITLE_CONNECT string
    TITLE_REDIRECT string
    MODE_READ string
    MODE_SAVE string
    MODE_TEST string
    MODE_REMV string
    MODE_MERG string
    MODE_DISTRIB string
    MODE_DECENTR string
    MODE_READ_MERG string
    MODE_SAVE_MERG string
    MODE_DISTRIB_READ string
    MODE_DISTRIB_SAVE string
    DEFAULT_HMAC_KEY []byte
    PACK_TIME time.Duration
    WAIT_TIME time.Duration
    ROUTE_NUM int
    BUFF_SIZE int
    PID_SIZE int
    SESSION_SIZE int
    MAXSIZE_ADDRESS int
    MAXSIZE_PACKAGE int
    CLIENT_NAME_SIZE int
    IS_NOTHING bool
    IS_DISTRIB bool
    IS_DECENTR bool
    HAS_CRYPTO bool
    HAS_ROUTING bool 
    HAS_FRIENDS bool
    CRYPTO_SPEED bool
    HANDLE_ROUTING bool
} {
    TEMPLATE: "0.0.0.0",
    CLIENT_NAME: "[CLIENT]",
    SEPARATOR: "\000\001\007\005\000",
    END_BYTES: "\000\005\007\001\000",
    NETWORK_NAME: "GENESIS",
    TITLE_CONNECT: "[TITLE:CONNECT]",
    TITLE_REDIRECT: "[TITLE:REDIRECT]",
    MODE_READ: "[MODE:READ]",
    MODE_SAVE: "[MODE:SAVE]",
    MODE_TEST: "[MODE:TEST]",
    MODE_REMV: "[MODE:REMV]",
    MODE_MERG: "[MODE:MERG]",
    MODE_DISTRIB: "[MODE:DISTRIB]",
    MODE_DECENTR: "[MODE:DECENTR]",
    MODE_READ_MERG: "[MODE:READ][MODE:MERG]",
    MODE_SAVE_MERG: "[MODE:SAVE][MODE:MERG]",
    MODE_DISTRIB_READ: "[MODE:DISTRIB][MODE:READ]",
    MODE_DISTRIB_SAVE: "[MODE:DISTRIB][MODE:SAVE]",
    DEFAULT_HMAC_KEY: []byte("DEFAULT-HMAC-KEY"),
    PACK_TIME: 10,
    WAIT_TIME: 5,
    ROUTE_NUM: 3,
    BUFF_SIZE: 1 << 9, // 512B
    PID_SIZE: 1 << 3, // 8B
    SESSION_SIZE: 1 << 5, // 32B
    MAXSIZE_ADDRESS: 1 << 5, // 32B
    MAXSIZE_PACKAGE: 1 << 10 << 10, // 1MiB
    CLIENT_NAME_SIZE: 1 << 4, // 16B
    IS_NOTHING: true,
    IS_DISTRIB: false,
    IS_DECENTR: false,
    HAS_CRYPTO: false,
    HAS_ROUTING: false,
    HAS_FRIENDS: false,
    CRYPTO_SPEED: false,
    HANDLE_ROUTING: false,
}

// Set up management.
func SettingsSet(settings SettingsType) []uint8 {
    var (
        list = make([]uint8, len(settings))
        i = 0
    )

    for name, data := range settings {
        switch data.(type) {
            case string: list[i] = stringSettings(name, data)
            case bool: list[i] = boolSettings(name, data)
            case int: list[i] = intSettings(name, data)
            default: list[i] = 2
        }
        i++
    }

    boolRelations()
    return list
}

// Get value by the key in setting.
func SettingsGet(key string) interface{} {
    switch key {
        case "TEMPLATE": return setting.TEMPLATE
        case "CLIENT_NAME": return setting.CLIENT_NAME
        case "SEPARATOR": return setting.SEPARATOR
        case "END_BYTES": return setting.END_BYTES
        case "NETWORK_NAME": return setting.NETWORK_NAME
        case "TITLE_CONNECT": return setting.TITLE_CONNECT
        case "TITLE_REDIRECT": return setting.TITLE_REDIRECT
        case "MODE_READ": return setting.MODE_READ
        case "MODE_SAVE": return setting.MODE_SAVE
        case "MODE_TEST": return setting.MODE_TEST
        case "MODE_REMV": return setting.MODE_REMV
        case "MODE_MERG": return setting.MODE_MERG
        case "MODE_DISTRIB": return setting.MODE_DISTRIB
        case "MODE_DECENTR": return setting.MODE_DECENTR
        case "MODE_READ_MERG": return setting.MODE_READ_MERG
        case "MODE_SAVE_MERG": return setting.MODE_SAVE_MERG
        case "MODE_DISTRIB_READ": return setting.MODE_DISTRIB_READ
        case "MODE_DISTRIB_SAVE": return setting.MODE_DISTRIB_SAVE
        case "DEFAULT_HMAC_KEY": return setting.DEFAULT_HMAC_KEY
        case "IS_DISTRIB": return setting.IS_DISTRIB
        case "IS_DECENTR": return setting.IS_DECENTR
        case "HAS_CRYPTO": return setting.HAS_CRYPTO
        case "HAS_ROUTING": return setting.HAS_ROUTING
        case "HAS_FRIENDS": return setting.HAS_FRIENDS
        case "CRYPTO_SPEED": return setting.CRYPTO_SPEED
        case "WAIT_TIME": return setting.WAIT_TIME
        case "PACK_TIME": return setting.PACK_TIME
        case "ROUTE_NUM": return setting.ROUTE_NUM
        case "BUFF_SIZE": return setting.BUFF_SIZE
        case "PID_SIZE": return setting.PID_SIZE
        case "SESSION_SIZE": return setting.SESSION_SIZE
        case "MAXSIZE_ADDRESS": return setting.MAXSIZE_ADDRESS
        case "MAXSIZE_PACKAGE": return setting.MAXSIZE_PACKAGE
        case "CLIENT_NAME_SIZE": return setting.CLIENT_NAME_SIZE
        case "HANDLE_ROUTING": return setting.HANDLE_ROUTING
    }
    return nil
}

func boolRelations() {
    if setting.IS_DISTRIB && setting.IS_DECENTR {
        setting.IS_DISTRIB = false
        setting.IS_DECENTR = false
        setting.IS_NOTHING = true
    } else if setting.IS_DISTRIB || setting.IS_DECENTR {
        setting.IS_NOTHING = false
    }

    if setting.CRYPTO_SPEED {
        setting.SESSION_SIZE = 1 << 4 // 16B
    } else {
        setting.SESSION_SIZE = 1 << 5 // 32B
    }
}

func stringSettings(name string, data interface{}) uint8 {
    result := data.(string)
    switch name {
        case "TEMPLATE": setting.TEMPLATE = result
        case "CLIENT_NAME": setting.CLIENT_NAME = result
        case "SEPARATOR": setting.SEPARATOR = result
        case "END_BYTES": setting.END_BYTES = result
        case "NETWORK_NAME": setting.NETWORK_NAME = result
        case "TITLE_CONNECT": setting.TITLE_CONNECT = result
        case "TITLE_REDIRECT": setting.TITLE_REDIRECT = result
        case "MODE_READ": setting.MODE_READ = result
        case "MODE_SAVE": setting.MODE_SAVE = result
        case "MODE_TEST": setting.MODE_TEST = result
        case "MODE_REMV": setting.MODE_REMV = result
        case "MODE_MERG": setting.MODE_MERG = result
        case "MODE_DISTRIB": setting.MODE_DISTRIB = result
        case "MODE_DECENTR": setting.MODE_DECENTR = result
        case "MODE_READ_MERG": setting.MODE_READ_MERG = result
        case "MODE_SAVE_MERG": setting.MODE_SAVE_MERG = result
        case "MODE_DISTRIB_READ": setting.MODE_DISTRIB_READ = result
        case "MODE_DISTRIB_SAVE": setting.MODE_DISTRIB_SAVE = result
        case "DEFAULT_HMAC_KEY": setting.DEFAULT_HMAC_KEY = []byte(result)
        default: return 1
    }
    return 0
}

func boolSettings(name string, data interface{}) uint8 {
    result := data.(bool)
    switch name {
        case "IS_DISTRIB": setting.IS_DISTRIB = result
        case "IS_DECENTR": setting.IS_DECENTR = result
        case "HAS_CRYPTO": setting.HAS_CRYPTO = result
        case "HAS_ROUTING": setting.HAS_ROUTING = result
        case "HAS_FRIENDS": setting.HAS_FRIENDS = result
        case "CRYPTO_SPEED": setting.CRYPTO_SPEED = result
        case "HANDLE_ROUTING": setting.HANDLE_ROUTING = result
        default: return 1
    }
    return 0
}

func intSettings(name string, data interface{}) uint8 {
    result := data.(int)
    switch name {
        case "WAIT_TIME": setting.WAIT_TIME = time.Duration(result)
        case "PACK_TIME": setting.PACK_TIME = time.Duration(result)
        case "ROUTE_NUM": setting.ROUTE_NUM = result
        case "BUFF_SIZE": setting.BUFF_SIZE = result
        case "PID_SIZE": setting.PID_SIZE = result
        case "SESSION_SIZE": setting.SESSION_SIZE = result
        case "MAXSIZE_ADDRESS": setting.MAXSIZE_ADDRESS = result
        case "MAXSIZE_PACKAGE": setting.MAXSIZE_PACKAGE = result
        case "CLIENT_NAME_SIZE": setting.CLIENT_NAME_SIZE = result
        default: return 1
    }
    return 0
}
