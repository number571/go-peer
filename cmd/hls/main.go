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
	// read config, database, key
	err := hlsDefaultInit()
	if err != nil {
		gLogger.Error(err.Error())
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
			gLogger.Warning(err.Error())
			continue
		}
		gLogger.Info(fmt.Sprintf("connected to '%s'", conn))
	}

	// if node in client mode
	// then run endless loop
	if gConfig.Address() == "" {
		gLogger.Info("Service is listening...")
		select {}
	}

	// run node in server mode
	gLogger.Info(fmt.Sprintf("Service is listening [%s]...", gConfig.Address()))
	err = node.Listen(gConfig.Address())
	if err != nil {
		gLogger.Error(err.Error())
		os.Exit(2)
	}
}

func routeHLS(node network.INode, msg local.IMessage) []byte {
	// load request from message's body
	requestBytes := msg.Body().Data()
	request := hlsnet.LoadRequest(requestBytes)
	if request == nil {
		return nil
	}

	// check already received data by hash
	hash := crypto.NewHasher(requestBytes).Bytes()
	if gDB.Exist(hash) {
		return nil
	}
	gDB.Push(hash)

	// get service's address by name
	info, ok := gConfig.GetService(request.Host())
	if !ok {
		return nil
	}

	// redirect bytes of request to another nodes
	if info.IsRedirect() {
		for _, recv := range gConfig.PubKeys() {
			go node.Request(
				local.NewRoute(recv, nil, nil),
				local.NewMessage([]byte(cPatternHLS), requestBytes),
			)
		}
	}

	// generate new request to serivce
	req, err := http.NewRequest(
		request.Method(),
		fmt.Sprintf("%s://%s%s", cProto, info.Address(), request.Path()),
		bytes.NewReader(request.Body()),
	)
	if err != nil {
		return nil
	}
	for key, val := range request.Head() {
		req.Header.Add(key, val)
	}

	// send request to service
	// and receive response from service
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	// send result to client
	return data
}
