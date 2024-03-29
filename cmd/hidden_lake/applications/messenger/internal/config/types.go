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
	UpdatePseudonym(string) error
	UpdateLanguage(language.ILanguage) error
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() string
}

type IConfigSettings interface {
	GetMessagesCapacity() uint64
	GetWorkSizeBits() uint64
	GetPseudonym() string
	GetStorageKey() string
	GetLanguage() language.ILanguage
}

type IAddress interface {
	GetInterface() string
	GetIncoming() string
	GetPPROF() string
}
