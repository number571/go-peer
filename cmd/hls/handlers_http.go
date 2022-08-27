package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	payload_adapter "github.com/number571/go-peer/modules/network/anonymity/adapters/payload"
)

func handleIndexHTTP(w http.ResponseWriter, r *http.Request) {
	response(w, hls_settings.CErrorNone, hls_settings.CTitlePattern)
}

func handleOnlineHTTP(w http.ResponseWriter, r *http.Request) {
	response(w, hls_settings.CErrorNone, strings.Join(gConnKeeper.InOnline(), ","))
}

func handleRequestHTTP(w http.ResponseWriter, r *http.Request) {
	var vRequest hls_settings.SRequest

	if r.Method != "POST" {
		response(w, hls_settings.CErrorMethod, "failed: incorrect method")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&vRequest)
	if err != nil {
		response(w, hls_settings.CErrorDecode, "failed: decode request")
		return
	}

	pubKey := asymmetric.LoadRSAPubKey(vRequest.Receiver)
	if pubKey == nil {
		response(w, hls_settings.CErrorPubKey, "failed: load public key")
		return
	}

	data := encoding.HexDecode(vRequest.HexData)
	if data == nil {
		response(w, hls_settings.CErrorPubKey, "failed: decode hex format data")
		return
	}

	resp, err := gNode.Request(
		pubKey,
		payload_adapter.NewPayload(hls_settings.CHeaderHLS, data),
	)
	if err != nil {
		response(w, hls_settings.CErrorResponse, "failed: response message")
		return
	}

	response(w, hls_settings.CErrorNone, encoding.HexEncode(resp))
}

func response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", hls_settings.CContentType)
	json.NewEncoder(w).Encode(&hls_settings.SResponse{
		Result: res,
		Return: ret,
	})
}
