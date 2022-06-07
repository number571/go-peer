package main

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func pageIndex(w http.ResponseWriter, r *http.Request) {
	response(w, hls_settings.CErrorNone, []byte(hls_settings.CTitlePattern))
}

func pageStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		response(w, hls_settings.CErrorMethod, []byte("failed: incorrect method"))
		return
	}

	var network []hls_settings.SStatusNetwork
	for _, info := range gNode.Checker().ListWithInfo() {
		network = append(network, hls_settings.SStatusNetwork{
			PubKey: info.PubKey().String(),
			Online: info.Online(),
		})
	}

	w.Header().Set("Content-Type", hls_settings.CContentType)
	json.NewEncoder(w).Encode(&hls_settings.SStatusResponse{
		PubKey:  gNode.Client().PubKey().String(),
		Network: network,
		SResponse: hls_settings.SResponse{
			Result: []byte("success"),
			Return: hls_settings.CErrorNone,
		},
	})
}

func pageMessage(w http.ResponseWriter, r *http.Request) {
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

	pubKey := crypto.LoadPubKey(vRequest.Receiver)
	if pubKey == nil {
		response(w, hls_settings.CErrorPubKey, []byte("failed: load public key"))
		return
	}

	inOnline := []crypto.IPubKey{}
	for _, val := range gNode.Checker().ListWithInfo() {
		if !val.Online() {
			continue
		}
		inOnline = append(inOnline, val.PubKey())
	}

	rand := crypto.NewPRNG()
	randSizeRoute := rand.Uint64() % hls_settings.CSizeRoute

	resp, err := gNode.Request(
		local.NewRoute(pubKey).
			WithRedirects(
				gPPrivKey,
				local.NewSelector(inOnline).
					Shuffle().
					Return(randSizeRoute),
			),
		local.NewMessage(vRequest.Title, vRequest.Data),
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
