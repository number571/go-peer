package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/internal/settings"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateLanguage(utils.ILanguage) error
}

type IConfig interface {
	IConfigSettingsHLM

	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() string
	GetStorageKey() string
	GetLanguage() utils.ILanguage
}

type IConfigSettingsHLM interface {
	IsValidHLM() bool
	settings.IConfigSettings

	GetKeySizeBits() uint64
	GetMessagesCapacity() uint64
}

type IAddress interface {
	GetInterface() string
	GetIncoming() string
}
