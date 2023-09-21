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
	service := network.NewNode(nodeSettings(serviceAddress))
	service.HandleFunc(serviceHeader, func(_ network.INode, c conn.IConn, reqBytes []byte) error {
		c.WritePayload(payload.NewPayload(
			serviceHeader,
			[]byte(fmt.Sprintf("echo: [%s]", string(reqBytes))),
		))
		return nil
	})

	if err := service.Run(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second) // wait

	conn, err := conn.NewConn(
		connSettings(),
		serviceAddress,
	)
	if err != nil {
		panic(err)
	}

	pld := payload.NewPayload(serviceHeader, []byte("hello, world!"))
	if err := conn.WritePayload(pld); err != nil {
		panic(err)
	}

	pld, err = conn.ReadPayload()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(pld.GetBody()))
}

func nodeSettings(serviceAddress string) network.ISettings {
	return network.NewSettings(&network.SSettings{
		FAddress:      serviceAddress,
		FCapacity:     (1 << 10),
		FMaxConnects:  1,
		FConnSettings: connSettings(),
		FWriteTimeout: time.Minute,
	})
}

func connSettings() conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FMessageSizeBytes: (1 << 10),
		FWaitReadDeadline: time.Hour,
		FReadDeadline:     time.Minute,
		FWriteDeadline:    time.Minute,
	})
}
