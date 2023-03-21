package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
)

func initConnKeeper(pCfg config.IConfig, pNode anonymity.INode) conn_keeper.IConnKeeper {
	return conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{
			FConnections: func() []string { return pCfg.GetConnections() },
			FDuration:    pNode.GetSettings().GetTimeWait(),
		}),
		pNode.GetNetworkNode(),
	)
}
