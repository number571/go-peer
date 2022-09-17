package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/database"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/settings"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

func handlePushHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response(w, hls_settings.CErrorMethod, "failed: incorrect method")
		return
	}

	msg, err := io.ReadAll(r.Body)
	if err != nil {
		response(w, hls_settings.CErrorResponse, "failed: response message")
		return
	}

	smsg := string(msg)

	pubKeyStr := r.Header[hls_settings.CHeaderPubKey][0]
	friend := asymmetric.LoadRSAPubKey(pubKeyStr)

	iam, err := gClient.PubKey()
	if err != nil {
		response(w, hls_settings.CErrorPubKey, "failed: get public key")
		return
	}

	err = gDB.Push(database.NewRelation(iam, friend), fmt.Sprintf("[%s]: %s", "friend", smsg))
	if err != nil {
		response(w, hls_settings.CErrorPubKey, "failed: push message to database")
		return
	}

	if gChannelPubKey != nil && pubKeyStr == gChannelPubKey.String() {
		fmt.Println(smsg)
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
