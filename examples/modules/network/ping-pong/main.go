package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/network/conn"
	"github.com/number571/go-peer/modules/payload"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

func main() {
	var (
		service1 = network.NewNode(network.NewSettings(&network.SSettings{}))
		service2 = network.NewNode(network.NewSettings(&network.SSettings{}))
	)

	service1.Handle(serviceHeader, handlerPayload("#1"))
	service2.Handle(serviceHeader, handlerPayload("#2"))

	go service1.Listen(serviceAddress)
	time.Sleep(time.Second) // wait

	_, err := service2.Connect(serviceAddress)
	if err != nil {
		panic(err)
	}

	service2.Broadcast(payload.NewPayload(
		serviceHeader,
		[]byte("0"),
	))

	select {}
}

func handlerPayload(serviceName string) network.IHandlerF {
	return func(n network.INode, c conn.IConn, p payload.IPayload) {
		time.Sleep(time.Second) // delay for view "ping-pong" game

		num, err := strconv.Atoi(string(p.Body()))
		if err != nil {
			panic(err)
		}

		val := "ping"
		if num%2 == 1 {
			val = "pong"
		}

		fmt.Printf("service '%s' got '%s#%d'\n", serviceName, val, num)
		n.Broadcast(payload.NewPayload(
			serviceHeader,
			[]byte(fmt.Sprintf("%d", num+1)),
		))
	}
}
