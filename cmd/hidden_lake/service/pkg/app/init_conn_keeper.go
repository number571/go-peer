package app

import (
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/connkeeper"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func (p *sApp) initConnKeeper(pNode network.INode) {
	p.fConnKeeper = connkeeper.NewConnKeeper(
		connkeeper.NewSettings(&connkeeper.SSettings{
			FConnections: func() []string { return p.fCfgW.GetConfig().GetConnections() },
			FDuration:    pkg_settings.CConnKeeperDuration,
		}),
		pNode,
	)
}
