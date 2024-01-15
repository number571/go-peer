package config

import (
	"github.com/number571/go-peer/internal/language"
	logger "github.com/number571/go-peer/internal/logger/std"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateLanguage(language.ILanguage) error
	UpdateSecretKeys(map[string]string) error
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetShare() bool
	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() string
	GetStorageKey() string
	GetLanguage() language.ILanguage
	GetSecretKeys() map[string]string
}

type IConfigSettings interface {
	GetMessagesCapacity() uint64
	GetWorkSizeBits() uint64
}

type IAddress interface {
	GetInterface() string
	GetIncoming() string
	GetPPROF() string
}
