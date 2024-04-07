package app

import (
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/connkeeper"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func (p *sApp) initConnKeeper(pNode network.INode) {
	p.fConnKeeper = connkeeper.NewConnKeeper(
		connkeeper.NewSettings(&connkeeper.SSettings{
			FConnections: func() []string { return p.fConfig.GetConnections() },
			FDuration:    hls_settings.CConnKeeperDuration,
		}),
		pNode,
	)
}
