package main

import (
	"flag"
	"fmt"

	"github.com/number571/go-peer/cmd/hls/app"
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/filesystem"
	"github.com/number571/go-peer/modules/friends"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/network/conn"
	"github.com/number571/go-peer/modules/queue"
	"github.com/number571/go-peer/modules/storage/database"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func initValues() error {
	var (
		inputKey string
	)

	flag.StringVar(&inputKey, "key", "priv.key", "input private key from file")
	flag.Parse()

	privKeyStr, err := filesystem.OpenFile(inputKey).Read()
	if err != nil {
		return err
	}

	privKey := asymmetric.LoadRSAPrivKey(string(privKeyStr))
	if privKey == nil {
		return fmt.Errorf("private key is invalid")
	}

	cfg, err := initConfig()
	if err != nil {
		return err
	}

	gApp = app.NewApp(cfg, initNode(cfg, privKey))
	return nil
}

func initConfig() (config.IConfig, error) {
	if filesystem.OpenFile(hls_settings.CPathCFG).IsExist() {
		return config.LoadConfig(hls_settings.CPathCFG)
	}
	initCfg := &config.SConfig{
		FAddress: &config.SAddress{
			FTCP:  "localhost:9571",
			FHTTP: "localhost:9572",
		},
	}
	return config.NewConfig(hls_settings.CPathCFG, initCfg)
}

func initNode(cfg config.IConfig, privKey asymmetric.IPrivKey) anonymity.INode {
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FRetryEnqueue: hls_settings.CRetryEnqueue,
			FTimeWait:     hls_settings.CWaitTime,
		}),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FPath:    hls_settings.CPathDB,
				FHashing: true,
			}),
		),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FCapacity:    hls_settings.CNetworkCapacity,
				FMaxConnects: hls_settings.CNetworkMaxConns,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:  cfg.Network(),
					FMessageSize: hls_settings.CMessageSize,
					FTimeWait:    hls_settings.CNetworkWaitTime,
				}),
			}),
		),
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     hls_settings.CQueueCapacity,
				FPullCapacity: hls_settings.CQueuePullCapacity,
				FDuration:     hls_settings.CQueueDuration,
			}),
			client.NewClient(
				client.NewSettings(&client.SSettings{
					FWorkSize:    hls_settings.CWorkSize,
					FMessageSize: hls_settings.CMessageSize,
				}),
				privKey,
			),
		),
		func() friends.IF2F {
			f2f := friends.NewF2F()
			for _, pubKey := range cfg.Friends() {
				f2f.Append(pubKey)
			}
			return f2f
		}(),
	)
}
