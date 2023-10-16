package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	logger "github.com/number571/go-peer/internal/logger/std"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateLanguage(utils.ILanguage) error
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() string
	GetStorageKey() string
	GetLanguage() utils.ILanguage
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
