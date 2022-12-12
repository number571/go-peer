package main

import (
	"fmt"
	"time"

	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/network/conn"
	"github.com/number571/go-peer/modules/payload"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

// client <-> service1 <-> service2
func main() {
	var (
		service1 = network.NewNode(network.NewSettings(&network.SSettings{}))
		service2 = network.NewNode(network.NewSettings(&network.SSettings{}))
	)

	service1.Handle(serviceHeader, handlerPingPong("#1"))
	service2.Handle(serviceHeader, handlerPingPong("#2"))

	go service1.Listen(serviceAddress)
	time.Sleep(time.Second) // wait

	_, err := service2.Connect(serviceAddress)
	if err != nil {
		panic(err)
	}

	conn, err := conn.NewConn(
		conn.NewSettings(&conn.SSettings{}),
		serviceAddress,
	)
	if err != nil {
		panic(err)
	}

	err = conn.Write(payload.NewPayload(
		serviceHeader,
		[]byte("hello, world!"),
	))
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
}

func handlerPingPong(serviceName string) network.IHandlerF {
	return func(n network.INode, c conn.IConn, reqBytes []byte) {
		defer n.Broadcast(payload.NewPayload(serviceHeader, reqBytes))
		fmt.Printf("service '%s' got '%s'\n", serviceName, string(reqBytes))
	}
}
