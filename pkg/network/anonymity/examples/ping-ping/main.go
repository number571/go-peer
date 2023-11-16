package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/queue_set"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
)

const (
	msgSize        = (100 << 10)
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

	service1.HandleFunc(serviceHeader, handler("#1"))
	service2.HandleFunc(serviceHeader, handler("#2"))

	service1.GetListPubKeys().AddPubKey(service2.GetMessageQueue().GetClient().GetPubKey())
	service2.GetListPubKeys().AddPubKey(service1.GetMessageQueue().GetClient().GetPubKey())

	if err := service1.Run(); err != nil {
		panic(err)
	}
	if err := service1.GetNetworkNode().Run(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	if err := service2.Run(); err != nil {
		panic(err)
	}
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

func newNode(serviceAddress, name, dbPath string) anonymity.INode {
	db, err := database.NewKeyValueDB(
		storage.NewSettings(&storage.SSettings{
			FPath: dbPath,
		}),
	)
	if err != nil {
		return nil
	}
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:   name,
			FRetryEnqueue:  0,
			FNetworkMask:   1,
			FFetchTimeWait: time.Minute,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{
				FInfo: os.Stdout,
				FWarn: os.Stdout,
				FErro: os.Stderr,
			}),
			func(arg logger.ILogArg) string {
				logBuilder, ok := arg.(anon_logger.ILogBuilder)
				if !ok {
					panic("got invalid log arg")
				}
				logGetter := logBuilder.Get()
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
				FDuration:     time.Second,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FWorkSizeBits:     10,
					FMessageSizeBytes: msgSize,
				}),
				asymmetric.NewRSAPrivKey(1024),
			),
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
		FWorkSizeBits:     10,
		FMessageSizeBytes: msgSize,
		FWaitReadDeadline: time.Hour,
		FReadDeadline:     time.Minute,
		FWriteDeadline:    time.Minute,
	})
}