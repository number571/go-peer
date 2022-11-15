package conn_keeper

import (
	"time"

	"github.com/number571/go-peer/modules"
)

type IConnKeeper interface {
	Settings() ISettings
	modules.IApp
}

type ISettings interface {
	GetConnections() []string
	GetDuration() time.Duration
}
