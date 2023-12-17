package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
	"github.com/number571/go-peer/pkg/utils"
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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	var (
		service1 = newNode(serviceAddress, "SERVICE-1", dbPath1)
		service2 = newNode("", "SERVICE-2", dbPath2)
	)

	defer func() {
		if err := closeNode(service1); err != nil {
			panic(err)
		}
		if err := closeNode(service2); err != nil {
			panic(err)
		}
	}()

	service1.HandleFunc(serviceHeader, handler("#1"))
	service2.HandleFunc(serviceHeader, handler("#2"))

	service1.GetListPubKeys().AddPubKey(service2.GetMessageQueue().GetClient().GetPubKey())
	service2.GetListPubKeys().AddPubKey(service1.GetMessageQueue().GetClient().GetPubKey())

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	go func() { _ = service1.Run(ctx1) }()
	if err := service1.GetNetworkNode().Listen(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	go func() { _ = service2.Run(ctx2) }()
	if err := service2.GetNetworkNode().AddConnection(serviceAddress); err != nil {
		panic(err)
	}

	err := service2.BroadcastPayload(
		service1.GetMessageQueue().GetClient().GetPubKey(),
		adapters.NewPayload(
			serviceHeader,
			[]byte("0"),
		),
	)
	if err != nil {
		panic(err)
	}

	<-shutdown
}

func handler(serviceName string) anonymity.IHandlerF {
	return func(node anonymity.INode, pubKey asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
		num, err := strconv.Atoi(string(reqBytes))
		if err != nil {
			panic(err)
		}

		val := "ping"
		if num%2 == 1 {
			val = "pong"
		}

		fmt.Printf("service '%s' got '%s#%d'\n", serviceName, val, num)

		err = node.BroadcastPayload(
			pubKey,
			adapters.NewPayload(
				serviceHeader,
				[]byte(fmt.Sprintf("%d", num+1)),
			),
		)
		if err != nil {
			panic(err)
		}

		return nil, nil
	}
}

func closeNode(node anonymity.INode) error {
	return utils.MergeErrors(
		node.GetWrapperDB().Close(),
		node.GetNetworkNode().Close(),
	)
}

func newNode(serviceAddress, name, dbPath string) anonymity.INode {
	db, err := database.NewKVDatabase(
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
				FDuration:     2 * time.Second,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FMessageSizeBytes: msgSize,
					FKeySizeBits:      keySize,
				}),
				asymmetric.NewRSAPrivKey(keySize),
			),
		).WithNetworkSettings(
			networkMask,
			net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits: workSize,
			}),
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
