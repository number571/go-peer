package main

import (
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

// client <-> service1 <-> service2
func main() {
	var (
		service1 = network.NewNode(nodeSettings(serviceAddress))
		service2 = network.NewNode(nodeSettings(""))
	)

	service1.HandleFunc(serviceHeader, handlerPingPong("#1"))
	service2.HandleFunc(serviceHeader, handlerPingPong("#2"))

	if err := service1.Run(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second) // wait

	if err := service2.AddConnection(serviceAddress); err != nil {
		panic(err)
	}

	conn, err := conn.NewConn(
		connSettings(),
		serviceAddress,
	)
	if err != nil {
		panic(err)
	}

	err = conn.WritePayload(payload.NewPayload(
		serviceHeader,
		[]byte("hello, world!"),
	))
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
}

func handlerPingPong(serviceName string) network.IHandlerF {
	return func(n network.INode, _ conn.IConn, reqBytes []byte) {
		defer n.BroadcastPayload(payload.NewPayload(serviceHeader, reqBytes))
		fmt.Printf("service '%s' got '%s'\n", serviceName, string(reqBytes))
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
		FMessageSizeBytes: (1 << 10),
		FFetchTimeWait:    1, // not used
	})
}
