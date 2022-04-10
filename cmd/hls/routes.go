package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	hlsnet "github.com/number571/go-peer/cmd/hls/network"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/network"
)

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
		for _, recv := range gConfig.F2F().Friends() {
			go node.Request(
				local.NewRoute(recv),
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
