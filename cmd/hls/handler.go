package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hls/config"
	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/payload"
)

func handleServiceTCP(cfg config.IConfig) anonymity.IHandlerF {
	return func(node anonymity.INode, sender asymmetric.IPubKey, pld payload.IPayload) []byte {
		// load request from message's body
		requestBytes := pld.Body()
		request := hls_network.LoadRequest(requestBytes)
		if request == nil {
			return nil
		}

		// get service's address by hostname
		address, ok := cfg.Service(request.Host())
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

		// set headers
		req.Header.Add(hls_settings.CHeaderPubKey, sender.String())
		for key, val := range request.Head() {
			if key == hls_settings.CHeaderPubKey {
				continue
			}
			req.Header.Add(key, val)
		}

		// send request to service
		// and receive response from service
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil
		}

		// send result to client
		return data
	}
}
