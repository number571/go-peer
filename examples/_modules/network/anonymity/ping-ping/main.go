package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/database"
)

const (
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
		service1 = newNode(serviceAddress, dbPath1)
		service2 = newNode("", dbPath2)
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

	if err := service2.GetNetworkNode().AddConnect(serviceAddress); err != nil {
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

	select {}
}

func handler(serviceName string) anonymity.IHandlerF {
	return func(node anonymity.INode, pubKey asymmetric.IPubKey, _, reqBytes []byte) []byte {
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
		return nil
	}
}

func newNode(serviceAddress, dbPath string) anonymity.INode {
	db, err := database.NewSQLiteDB(
		database.NewSettings(&database.SSettings{FPath: dbPath}),
	)
	if err != nil {
		return nil
	}
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		anonymity.NewWrapperDB().Set(db),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FConnSettings: conn.NewSettings(&conn.SSettings{}),
				FAddress:      serviceAddress,
			}),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{}),
			client.NewClient(
				message.NewSettings(&message.SSettings{}),
				asymmetric.NewRSAPrivKey(1024),
			),
		),
		asymmetric.NewListPubKeys(),
	)
}
