// HLS - Hidden Lake Service
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	hlsnet "github.com/number571/go-peer/cmd/hls/network"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/network"
)

func main() {
	err := hlsDefaultInit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create node from client
	node := network.NewNode(gClient).
		Handle([]byte(cPatternHLS), routeHLS)

	// turn on f2f mode
	node.F2F().Set(gConfig.F2F())
	for _, pubKey := range gConfig.PubKeys() {
		node.F2F().Append(pubKey)
	}

	// connect to open nodes
	for _, conn := range gConfig.Connections() {
		err := node.Connect(conn)
		if err != nil {
			fmt.Println(err)
		}
	}

	if gConfig.Address() == "" {
		select {}
	}

	err = node.Listen(gConfig.Address())
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func routeHLS(node network.INode, msg local.IMessage) []byte {
	request := hlsnet.LoadRequest(msg.Body().Data())
	if request == nil {
		return nil
	}

	hash := crypto.NewHasher(request.Body()).Bytes()
	if gDB.Exist(hash) {
		return nil
	}
	gDB.Push(hash)

	info, ok := gConfig.GetService(request.Host())
	if !ok {
		return nil
	}

	if info.IsRedirect() {
		for _, recv := range gConfig.PubKeys() {
			go node.Request(
				local.NewRoute(recv, nil, nil),
				local.NewMessage([]byte(cPatternHLS), request.Body()),
			)
		}
	}

	req, err := http.NewRequest(
		request.Method(),
		info.Address()+request.Path(),
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
