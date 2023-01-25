package conn_keeper

import (
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/types"
)

type IConnKeeper interface {
	Settings() ISettings
	Network() network.INode
	types.IApp
}

type ISettings interface {
	GetConnections() []string
	GetDuration() time.Duration
}
