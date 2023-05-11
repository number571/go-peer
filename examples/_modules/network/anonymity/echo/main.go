package main

import (
	"fmt"
	"os"
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
		service = newNode(serviceAddress, dbPath1)
		client  = newNode("", dbPath2)
	)

	service.HandleFunc(serviceHeader, func(_ anonymity.INode, _ asymmetric.IPubKey, _, reqBytes []byte) []byte {
		return []byte(fmt.Sprintf("echo: [%s]", string(reqBytes)))
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

	if err := client.GetNetworkNode().AddConnect(serviceAddress); err != nil {
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
				FAddress:      serviceAddress,
				FConnSettings: conn.NewSettings(&conn.SSettings{}),
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
