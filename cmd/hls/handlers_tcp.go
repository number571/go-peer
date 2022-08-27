package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	hlsnet "github.com/number571/go-peer/cmd/hls/network"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/payload"
)

func handleTCP(node anonymity.INode, _ asymmetric.IPubKey, pld payload.IPayload) []byte {
	// load request from message's body
	requestBytes := pld.Body()
	request := hlsnet.LoadRequest(requestBytes)
	if request == nil {
		return nil
	}

	// get service's address by hostname
	address, ok := gConfig.Service(request.Host())
	if !ok {
		return nil
	}

	// generate new request to serivce
	req, err := http.NewRequest(
		request.Method(),
		fmt.Sprintf("http://%s%s", address, request.Path()),
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
