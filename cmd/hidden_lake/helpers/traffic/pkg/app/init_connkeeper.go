package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/config"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/connkeeper"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func initConnKeeper(pCfg config.IConfig, pNode network.INode) connkeeper.IConnKeeper {
	return connkeeper.NewConnKeeper(
		connkeeper.NewSettings(&connkeeper.SSettings{
			FConnections: func() []string { return pCfg.GetConnections() },
			FDuration:    hls_settings.CConnKeeperDuration,
		}),
		pNode,
	)
}
