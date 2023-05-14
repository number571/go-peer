package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

func main() {
	var (
		service1 = network.NewNode(nodeSettings(serviceAddress))
		service2 = network.NewNode(nodeSettings(""))
	)

	service1.HandleFunc(serviceHeader, handler("#1"))
	service2.HandleFunc(serviceHeader, handler("#2"))

	if err := service1.Run(); err != nil {
		panic(err)
	}
	defer service1.Stop()

	time.Sleep(time.Second) // wait

	if err := service2.AddConnect(serviceAddress); err != nil {
		panic(err)
	}

	service2.BroadcastPayload(payload.NewPayload(
		serviceHeader,
		[]byte("0"),
	))

	select {}
}

func handler(serviceName string) network.IHandlerF {
	return func(n network.INode, c conn.IConn, reqBytes []byte) {
		time.Sleep(time.Second) // delay for view "ping-pong" game

		num, err := strconv.Atoi(string(reqBytes))
		if err != nil {
			panic(err)
		}

		val := "ping"
		if num%2 == 1 {
			val = "pong"
		}

		fmt.Printf("service '%s' got '%s#%d'\n", serviceName, val, num)
		n.BroadcastPayload(payload.NewPayload(
			serviceHeader,
			[]byte(fmt.Sprintf("%d", num+1)),
		))
	}
}

func nodeSettings(serviceAddress string) network.ISettings {
	return network.NewSettings(&network.SSettings{
		FAddress:      serviceAddress,
		FCapacity:     (1 << 10),
		FMaxConnects:  1,
		FConnSettings: connSettings(),
	})
}

func connSettings() conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FMessageSize:   (1 << 10),
		FLimitVoidSize: 1, // not used
		FFetchTimeWait: 1, // not used
	})
}
