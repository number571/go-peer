package main

import (
	"fmt"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/friends"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
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
		service = newNode(dbPath1)
		client  = newNode(dbPath2)
	)

	service.Handle(serviceHeader, func(_ anonymity.INode, _ asymmetric.IPubKey, _, reqBytes []byte) []byte {
		return []byte(fmt.Sprintf("echo: [%s]", string(reqBytes)))
	})

	service.F2F().Append(client.Queue().Client().PubKey())
	client.F2F().Append(service.Queue().Client().PubKey())

	if err := service.Run(); err != nil {
		panic(err)
	}

	if err := client.Run(); err != nil {
		panic(err)
	}

	go service.Network().Listen(serviceAddress)
	time.Sleep(time.Second)

	if _, err := client.Network().Connect(serviceAddress); err != nil {
		panic(err)
	}

	res, err := client.Request(
		service.Queue().Client().PubKey(),
		anonymity.NewPayload(serviceHeader, []byte("hello, world!")),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}

func newNode(dbPath string) anonymity.INode {
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{FPath: dbPath}),
		),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FConnSettings: conn.NewSettings(&conn.SSettings{}),
			}),
		),
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{}),
			client.NewClient(
				client.NewSettings(&client.SSettings{}),
				asymmetric.NewRSAPrivKey(1024),
			),
		),
		friends.NewF2F(),
	)
}
