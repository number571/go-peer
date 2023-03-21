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
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleServiceTCP(pCfg config.IConfig) anonymity.IHandlerF {
	return func(_ anonymity.INode, sender asymmetric.IPubKey, msgHash, reqBytes []byte) []byte {
		// load request from message's body
		loadReq, err := request.LoadRequest(reqBytes)
		if err != nil {
			return nil
		}

		// get service's address by hostname
		address, ok := pCfg.GetService(loadReq.Host())
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

		// set service headers
		pushReq.Header.Add(pkg_settings.CHeaderPubKey, sender.ToString())
		pushReq.Header.Add(pkg_settings.CHeaderMsgHash, encoding.HexEncode(msgHash))

		// append headers from request
		for key, val := range loadReq.Head() {
			switch key {
			case pkg_settings.CHeaderPubKey, pkg_settings.CHeaderMsgHash:
				continue
			default:
				pushReq.Header.Add(key, val)
			}
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
