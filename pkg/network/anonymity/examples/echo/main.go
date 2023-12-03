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
	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/queue_set"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
)

const (
	workSize       = 10
	keySize        = 1024
	msgSize        = (8 << 10)
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

const (
	dbPath1 = "database1.db"
	dbPath2 = "database2.db"
)

func deleteDBs() {
	os.RemoveAll(dbPath1)
	os.RemoveAll(dbPath2)
}

func main() {
	deleteDBs()
	defer deleteDBs()

	var (
		service = newNode(serviceAddress, "SERVICE-1", dbPath1)
		client  = newNode("", "SERVICE-2", dbPath2)
	)

	service.HandleFunc(serviceHeader, func(_ anonymity.INode, _ asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
		return []byte(fmt.Sprintf("echo: [%s]", string(reqBytes))), nil
	})

	service.GetListPubKeys().AddPubKey(client.GetMessageQueue().GetClient().GetPubKey())
	client.GetListPubKeys().AddPubKey(service.GetMessageQueue().GetClient().GetPubKey())

	if err := service.Run(); err != nil {
		panic(err)
	}
	if err := service.GetNetworkNode().Run(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second) // wait

	if err := client.Run(); err != nil {
		panic(err)
	}
	if err := client.GetNetworkNode().AddConnection(serviceAddress); err != nil {
		panic(err)
	}

	res, err := client.FetchPayload(
		service.GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(serviceHeader, []byte("hello, world!")),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}

func newNode(serviceAddress, name, dbPath string) anonymity.INode {
	db, err := database.NewKeyValueDB(
		storage.NewSettings(&storage.SSettings{
			FPath: dbPath,
		}),
	)
	if err != nil {
		return nil
	}
	networkMask := uint64(1)
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:   name,
			FRetryEnqueue:  0,
			FNetworkMask:   networkMask,
			FFetchTimeWait: time.Minute,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{
				FInfo: os.Stdout,
				FWarn: os.Stdout,
				FErro: os.Stderr,
			}),
			func(arg logger.ILogArg) string {
				logGetterFactory, ok := arg.(anon_logger.ILogGetterFactory)
				if !ok {
					panic("got invalid log arg")
				}
				logGetter := logGetterFactory.Get()
				return fmt.Sprintf(
					"%s|%02xT|%XH|%dP|%dB",
					logGetter.GetService(),
					logGetter.GetType(),
					logGetter.GetHash(),
					logGetter.GetProof(),
					logGetter.GetSize(),
				)
			},
		),
		anonymity.NewWrapperDB().Set(db),
		network.NewNode(
			nodeSettings(serviceAddress),
			queue_set.NewQueueSet(
				queue_set.NewSettings(&queue_set.SSettings{
					FCapacity: (1 << 10),
				}),
			),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FMainCapacity: (1 << 4),
				FPoolCapacity: (1 << 4),
				FDuration:     5 * time.Second,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FMessageSizeBytes: msgSize,
					FKeySizeBits:      keySize,
				}),
				asymmetric.NewRSAPrivKey(keySize),
			),
			func() (uint64, net_message.ISettings) {
				return networkMask, net_message.NewSettings(&net_message.SSettings{
					FWorkSizeBits: workSize,
				})
			},
		),
		asymmetric.NewListPubKeys(),
	)
}

func nodeSettings(serviceAddress string) network.ISettings {
	return network.NewSettings(&network.SSettings{
		FAddress:      serviceAddress,
		FMaxConnects:  1,
		FConnSettings: connSettings(),
		FWriteTimeout: time.Minute,
		FReadTimeout:  time.Minute,
	})
}

func connSettings() conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FWorkSizeBits:     workSize,
		FMessageSizeBytes: msgSize,
		FWaitReadDeadline: time.Hour,
		FReadDeadline:     time.Minute,
		FWriteDeadline:    time.Minute,
	})
}
