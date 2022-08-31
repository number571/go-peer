package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	hlm_settings "github.com/number571/go-peer/cmd/hlm/settings"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func handlePushHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response(w, hls_settings.CErrorMethod, "failed: incorrect method")
		return
	}

	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response(w, hls_settings.CErrorResponse, "failed: response message")
		return
	}

	if gChannelPubKey != nil && r.Header[hls_settings.CHeaderPubKey][0] == gChannelPubKey.String() {
		fmt.Println(string(msg))
	}

	response(w, hls_settings.CErrorNone, hlm_settings.CTitlePattern)
}

func response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", hls_settings.CContentType)
	json.NewEncoder(w).Encode(&hls_settings.SResponse{
		FResult: res,
		FReturn: ret,
	})
}
