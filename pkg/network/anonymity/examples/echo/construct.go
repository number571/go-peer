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
		newKVDatabase(serviceName),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      address,
				FMaxConnects:  256,
				FReadTimeout:  time.Minute,
				FWriteTimeout: time.Minute,
				FConnSettings: newConnSettings(workSize, msgSize),
			}),
			newVSettings(networkKey),
			lru.NewLRUCache(1024),
		),
		queue.NewMessageQueueProcessor(
			queue.NewSettings(&queue.SSettings{
				FNetworkMask:      networkMask,
				FWorkSizeBits:     workSize,
				FQueuePeriod:      2 * time.Second,
				FMainPoolCapacity: 32,
				FRandPoolCapacity: 32,
			}),
			queue.NewVSettings(&queue.SVSettings{
				FNetworkKey: networkKey,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      keySize,
					FMessageSizeBytes: msgSize,
				}),
				asymmetric.NewRSAPrivKey(keySize),
			),
		),
		asymmetric.NewListPubKeys(),
	)
}

func newKVDatabase(serviceName string) database.IKVDatabase {
	db, err := database.NewKVDatabase(
		database.NewSettings(&database.SSettings{
			FPath: "./database_" + serviceName + ".db",
		}),
	)
	if err != nil {
		panic(err)
	}
	return db
}

func newVSettings(nKey string) conn.IVSettings {
	return conn.NewVSettings(&conn.SVSettings{
		FNetworkKey: nKey,
	})
}

func newConnSettings(wSize, mSize uint64) conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FLimitMessageSizeBytes: mSize + 4096,
		FWorkSizeBits:          wSize,
		FWaitReadTimeout:       time.Hour,
		FDialTimeout:           time.Minute,
		FReadTimeout:           time.Minute,
		FWriteTimeout:          time.Minute,
	})
}
