package main

import (
	"context"
	"fmt"
	"os"
	"time"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/anonymity/qb/adapters"
	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
)

const (
	networkKey  = "network-key"
	networkMask = uint32(0x11223344)
	keySize     = uint64(1024)
	msgSize     = uint64(8192)
	workSize    = uint64(10)
)

func newNode(serviceName, address string) (network.INode, anonymity.INode) {
	msgChan := make(chan layer1.IMessage)
	networkNode := network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      address,
			FMaxConnects:  256,
			FReadTimeout:  time.Minute,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSettings: layer1.NewSettings(&layer1.SSettings{
					FWorkSizeBits: workSize,
				}),
				FLimitMessageSizeBytes: msgSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		}),
		cache.NewLRUCache(1024),
	).HandleFunc(
		networkMask,
		func(ctx context.Context, _ network.INode, _ conn.IConn, msg layer1.IMessage) error {
			msgChan <- msg
			return nil
		},
	)
	anonymityNode := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:  serviceName,
			FFetchTimeout: time.Minute,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{
				FInfo: os.Stdout,
				FWarn: os.Stdout,
				FErro: os.Stderr,
			}),
			func(ia logger.ILogArg) string {
				logGetter, ok := ia.(anon_logger.ILogGetter)
				if !ok {
					panic("got invalid log arg")
				}
				return fmt.Sprintf(
					"name=%s code=%02x hash=%X proof=%08d bytes=%d",
					logGetter.GetService(),
					logGetter.GetType(),
					logGetter.GetHash()[:16],
					logGetter.GetProof(),
					logGetter.GetSize(),
				)
			},
		),
		adapters.NewAdapterByFuncs(
			func(ctx context.Context, msg layer1.IMessage) error {
				return networkNode.BroadcastMessage(ctx, msg)
			},
			func(ctx context.Context) (layer1.IMessage, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case msg := <-msgChan:
					_ = networkNode.BroadcastMessage(ctx, msg)
					return msg, nil
				}
			},
		),
		func() database.IKVDatabase {
			db, err := database.NewKVDatabase("./database_" + serviceName + ".db")
			if err != nil {
				panic(err)
			}
			return db
		}(),
		queue.NewQBProblemProcessor(
			queue.NewSettings(&queue.SSettings{
				FMessageConstructSettings: layer1.NewConstructSettings(&layer1.SConstructSettings{
					FSettings: layer1.NewSettings(&layer1.SSettings{
						FWorkSizeBits: workSize,
					}),
				}),
				FNetworkMask:  networkMask,
				FQueuePeriod:  2 * time.Second,
				FConsumersCap: 1,
				FQueuePoolCap: [2]uint64{32, 32},
			}),
			client.NewClient(
				asymmetric.NewPrivKey(),
				msgSize,
			),
		),
	)
	return networkNode, anonymityNode
}
