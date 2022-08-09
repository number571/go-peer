package main

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/payload"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func pageIndex(w http.ResponseWriter, r *http.Request) {
	response(w, hls_settings.CErrorNone, []byte(hls_settings.CTitlePattern))
}

func pageSend(w http.ResponseWriter, r *http.Request) {
	var vRequest hls_settings.SRequest

	if r.Method != "POST" {
		response(w, hls_settings.CErrorMethod, []byte("failed: incorrect method"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&vRequest)
	if err != nil {
		response(w, hls_settings.CErrorDecode, []byte("failed: decode request"))
		return
	}

	pubKey := asymmetric.LoadRSAPubKey(vRequest.Receiver)
	if pubKey == nil {
		response(w, hls_settings.CErrorPubKey, []byte("failed: load public key"))
		return
	}

	resp, err := gNode.Request(
		pubKey,
		payload.NewPayload(uint64(hls_settings.CHeaderHLS), vRequest.Data),
	)
	if err != nil {
		response(w, hls_settings.CErrorResponse, []byte("failed: response message"))
		return
	}

	response(w, hls_settings.CErrorNone, resp)
}

func response(w http.ResponseWriter, ret int, res []byte) {
	w.Header().Set("Content-Type", hls_settings.CContentType)
	json.NewEncoder(w).Encode(&hls_settings.SResponse{
		Result: res,
		Return: ret,
	})
}
