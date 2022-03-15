package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	hlsnet "github.com/number571/go-peer/cmd/hls/network"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/network"
)

func main() {
	hlsDefaultInit()
	fmt.Println("Service is listening...")

	node := network.NewNode(gClient).
		Handle([]byte(cPatternHLS), routeHLS)

	for _, conn := range gConfig.Connections() {
		err := node.Connect(conn)
		if err != nil {
			fmt.Println(err)
		}
	}

	if gConfig.Address() == "" {
		select {}
	}

	err := node.Listen(gConfig.Address())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func routeHLS(client local.IClient, msg local.IMessage) []byte {
	request := hlsnet.LoadRequest(msg.Body().Data())
	if request == nil {
		return nil
	}

	addr, ok := gConfig.GetService(request.Host())
	if !ok {
		return nil
	}

	req, err := http.NewRequest(
		request.Method(),
		addr+request.Path(),
		bytes.NewReader(request.Body()),
	)
	if err != nil {
		return nil
	}

	for key, val := range request.Head() {
		req.Header.Add(key, val)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return data
}
