package conn_keeper

import (
	"time"

	"github.com/number571/go-peer/pkg/types"
)

type IConnKeeper interface {
	Settings() ISettings
	types.IApp
}

type ISettings interface {
	GetConnections() []string
	GetDuration() time.Duration
}
