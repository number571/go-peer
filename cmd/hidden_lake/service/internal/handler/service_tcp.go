package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleServiceTCP(cfg config.IConfig) anonymity.IHandlerF {
	return func(node anonymity.INode, sender asymmetric.IPubKey, reqBytes []byte) []byte {
		// load request from message's body
		loadReq := request.LoadRequest(reqBytes)
		if loadReq == nil {
			return nil
		}

		// get service's address by hostname
		address, ok := cfg.Service(loadReq.Host())
		if !ok {
			return nil
		}

		// generate new request to serivce
		pushReq, err := http.NewRequest(
			loadReq.Method(),
			fmt.Sprintf("http://%s%s", address, loadReq.Path()),
			bytes.NewReader(loadReq.Body()),
		)
		if err != nil {
			return nil
		}

		// set headers
		pushReq.Header.Add(pkg_settings.CHeaderPubKey, sender.String())
		for key, val := range loadReq.Head() {
			if key == pkg_settings.CHeaderPubKey {
				continue
			}
			pushReq.Header.Add(key, val)
		}

		// send request to service
		// and receive response from service
		resp, err := http.DefaultClient.Do(pushReq)
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
