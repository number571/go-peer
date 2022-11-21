package app

import (
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/network/conn_keeper"
)

func initConnKeeper(cfg config.IConfig, node anonymity.INode) conn_keeper.IConnKeeper {
	return conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{
			FConnections: func() []string { return cfg.Connections() },
			FDuration:    node.Settings().GetTimeWait(),
		}),
		node.Network(),
	)
}
