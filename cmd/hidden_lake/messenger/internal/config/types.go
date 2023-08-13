package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/internal/settings"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateBackupConnections([]string) error
	UpdateLanguage(utils.ILanguage) error
}

type IConfig interface {
	IConfigSettingsHLM

	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() string
	GetStorageKey() string
	GetBackupConnections() []string
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
