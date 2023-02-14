package main

import (
	"fmt"
	"os"
	"strconv"
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
		service1 = newNode(dbPath1)
		service2 = newNode(dbPath2)
	)

	service1.Handle(serviceHeader, handler("#1"))
	service2.Handle(serviceHeader, handler("#2"))

	service1.F2F().Append(service2.Queue().Client().PubKey())
	service2.F2F().Append(service1.Queue().Client().PubKey())

	if err := service1.Run(); err != nil {
		panic(err)
	}

	if err := service2.Run(); err != nil {
		panic(err)
	}

	go service1.Network().Listen(serviceAddress)
	time.Sleep(time.Second)

	if _, err := service2.Network().Connect(serviceAddress); err != nil {
		panic(err)
	}

	err := service2.Broadcast(
		service1.Queue().Client().PubKey(),
		anonymity.NewPayload(
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

		err = node.Broadcast(
			pubKey,
			anonymity.NewPayload(
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
