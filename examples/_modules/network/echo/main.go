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

func main() {
	service := network.NewNode(network.NewSettings(&network.SSettings{FAddress: serviceAddress}))
	service.HandleFunc(serviceHeader, func(n network.INode, c conn.IConn, reqBytes []byte) {
		c.WritePayload(payload.NewPayload(
			serviceHeader,
			[]byte(fmt.Sprintf("echo: [%s]", string(reqBytes))),
		))
	})

	if err := service.Run(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second) // wait

	conn, err := conn.NewConn(
		conn.NewSettings(&conn.SSettings{}),
		serviceAddress,
	)
	if err != nil {
		panic(err)
	}

	pld, err := conn.FetchPayload(payload.NewPayload(
		serviceHeader,
		[]byte("hello, world!")),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(pld.GetBody()))
}
