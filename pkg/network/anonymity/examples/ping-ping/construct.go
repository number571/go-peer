package main

import (
	"fmt"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/cache/lru"
	"github.com/number571/go-peer/pkg/storage/database"
)

const (
	networkKey  = "network-key"
	networkMask = uint32(0x11223344)
	keySize     = uint64(1024)
	msgSize     = uint64(8192)
	workSize    = uint64(10)
)

func newNode(serviceName, address string) anonymity.INode {
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:  serviceName,
			FNetworkMask:  networkMask,
			FFetchTimeout: time.Minute,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{
				FInfo: os.Stdout,
				FWarn: os.Stdout,
				FErro: os.Stdout,
			}),
			func(ia logger.ILogArg) string {
				logGetterFactory, ok := ia.(anon_logger.ILogGetterFactory)
				if !ok {
					panic("got invalid log arg")
				}

				logGetter := logGetterFactory.Get()
				if logGetter.GetHash() == nil {
					// request with null hash from sender
					return ""
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
		func() database.IKVDatabase {
			db, err := database.NewKVDatabase(
				database.NewSettings(&database.SSettings{
					FPath: "./database_" + serviceName + ".db",
				}),
			)
			if err != nil {
				panic(err)
			}
			return db
		}(),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      address,
				FMaxConnects:  256,
				FReadTimeout:  time.Minute,
				FWriteTimeout: time.Minute,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FMessageSettings: net_message.NewSettings(&net_message.SSettings{
						FWorkSizeBits: workSize,
					}),
					FLimitMessageSizeBytes: msgSize,
					FWaitReadTimeout:       time.Hour,
					FDialTimeout:           time.Minute,
					FReadTimeout:           time.Minute,
					FWriteTimeout:          time.Minute,
				}),
			}),
			lru.NewLRUCache(1024),
		),
		queue.NewQBProblemProcessor(
			queue.NewSettings(&queue.SSettings{
				FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
					FSettings: net_message.NewSettings(&net_message.SSettings{
						FWorkSizeBits: workSize,
					}),
				}),
				FNetworkMask:      networkMask,
				FQueuePeriod:      2 * time.Second,
				FMainPoolCapacity: 32,
				FRandPoolCapacity: 32,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      keySize,
					FMessageSizeBytes: msgSize,
				}),
				asymmetric.NewRSAPrivKey(keySize),
			),
			asymmetric.NewRSAPrivKey(keySize).GetPubKey(),
		),
		asymmetric.NewListPubKeys(),
	)
}
