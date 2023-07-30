package handler

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
)

func HandleNetworkRequestAPI(pWrapper config.IWrapper, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		var vPush pkg_settings.SRequest

		if pR.Method != http.MethodPost && pR.Method != http.MethodPut {
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vPush); err != nil {
			api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		friends := pWrapper.GetConfig().GetFriends()
		pubKey, ok := friends[vPush.FReceiver]
		if !ok {
			api.Response(pW, http.StatusBadRequest, "failed: load public key")
			return
		}

		data := encoding.HexDecode(vPush.FHexData)
		if data == nil {
			api.Response(pW, http.StatusTeapot, "failed: decode hex format data")
			return
		}

		if _, err := request.LoadRequest(data); err != nil {
			api.Response(pW, http.StatusForbidden, "failed: decode request")
			return
		}

		switch pR.Method {
		case http.MethodPut:
			err := pNode.BroadcastPayload(
				pubKey,
				adapters.NewPayload(pkg_settings.CServiceMask, data),
			)
			if err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: broadcast message")
				return
			}
			api.Response(pW, http.StatusOK, "success: broadcast")
			return
		case http.MethodPost:
			respBytes, err := pNode.FetchPayload(
				pubKey,
				adapters.NewPayload(pkg_settings.CServiceMask, data),
			)
			if err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: get response bytes")
				return
			}
			if _, err := response.LoadResponse(respBytes); err != nil {
				api.Response(pW, http.StatusNotExtended, "failed: load response bytes")
				return
			}
			api.Response(pW, http.StatusOK, string(respBytes))
			return
		}
	}
}
