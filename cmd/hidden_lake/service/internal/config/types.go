package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateNetworkKey(string) error
	UpdateConnections([]string) error
	UpdateBackupConnections([]string) error
	UpdateFriends(map[string]asymmetric.IPubKey) error
}

type IConfigSettings interface {
	message.ISettings

	GetKeySizeBits() uint64
	GetQueuePeriodMS() uint64
	GetLimitVoidSizeBytes() uint64
	GetMessagesCapacity() uint64
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetNetworkKey() string
	GetConnections() []string
	GetBackupConnections() []string
	GetFriends() map[string]asymmetric.IPubKey
	GetService(string) (string, bool)
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
	GetPPROF() string
}
