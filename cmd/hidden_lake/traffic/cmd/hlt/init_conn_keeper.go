package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/settings"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
)

func initConnKeeper(cfg config.IConfig, db database.IKeyValueDB) conn_keeper.IConnKeeper {
	return conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{
			FConnections: func() []string { return []string{cfg.Connection()} },
			FDuration:    hlt_settings.CNetworkWaitTime,
		}),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FMaxConnects: 1, // only to HLS
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:  cfg.Network(),
					FMessageSize: hlt_settings.CMessageSize,
					FTimeWait:    hlt_settings.CNetworkWaitTime,
				}),
			}),
		).Handle(
			hlt_settings.CNetworkMask,
			func(_ network.INode, _ conn.IConn, reqBytes []byte) {
				msg := message.LoadMessage(
					reqBytes,
					db.Settings().GetMessageSize(),
					db.Settings().GetWorkSize(),
				)
				if msg == nil {
					// TODO: log
					return
				}
				strHash := encoding.HexEncode(msg.Body().Hash())
				if _, err := db.Load(strHash); err == nil {
					// TODO: log
					return
				}
				if err := db.Push(msg); err != nil {
					// TODO: log
					return
				}
			},
		),
	)
}
