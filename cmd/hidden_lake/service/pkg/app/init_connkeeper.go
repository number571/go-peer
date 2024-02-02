package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/connkeeper"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func initConnKeeper(pCfg config.IConfig, pNode anonymity.INode) connkeeper.IConnKeeper {
	return connkeeper.NewConnKeeper(
		connkeeper.NewSettings(&connkeeper.SSettings{
			FConnections: func() []string { return pCfg.GetConnections() },
			FDuration:    pkg_settings.CConnKeeperDuration,
		}),
		pNode.GetNetworkNode(),
	)
}
