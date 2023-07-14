package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn_keeper"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func initConnKeeper(pCfg config.IConfig, pNode network.INode) conn_keeper.IConnKeeper {
	return conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{
			FConnections: func() []string { return pCfg.GetConnections() },
			FDuration:    hls_settings.CConnKeeperDuration,
		}),
		pNode,
	)
}
