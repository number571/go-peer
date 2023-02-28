package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
)

func initConnKeeper(cfg config.IConfig, node anonymity.INode) conn_keeper.IConnKeeper {
	return conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{
			FConnections: func() []string { return cfg.Connections() },
			FDuration:    node.GetSettings().GetTimeWait(),
		}),
		node.GetNetworkNode(),
	)
}
