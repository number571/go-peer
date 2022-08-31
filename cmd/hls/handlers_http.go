package main

import (
	"encoding/json"
	"net/http"
	"sort"
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
	conns := gNode.Network().Connections()
	inOnline := make([]string, 0, len(conns))
	for _, conn := range conns {
		inOnline = append(inOnline, conn.Socket().RemoteAddr().String())
	}
	sort.SliceStable(inOnline, func(i, j int) bool {
		return inOnline[i] < inOnline[j]
	})
	response(w, hls_settings.CErrorNone, strings.Join(inOnline, ","))
}

func handlePubKeyHTTP(w http.ResponseWriter, r *http.Request) {
	response(w, hls_settings.CErrorNone, gNode.Queue().Client().PubKey().String())
}

func handleFriendsHTTP(w http.ResponseWriter, r *http.Request) {
	listPubKeys := gNode.F2F().List()
	listPubKeysStr := make([]string, 0, len(listPubKeys))
	for _, pubKey := range listPubKeys {
		listPubKeysStr = append(listPubKeysStr, pubKey.String())
	}
	response(w, hls_settings.CErrorNone, strings.Join(listPubKeysStr, ","))
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

	pubKey := asymmetric.LoadRSAPubKey(vRequest.FReceiver)
	if pubKey == nil {
		response(w, hls_settings.CErrorPubKey, "failed: load public key")
		return
	}

	data := encoding.HexDecode(vRequest.FHexData)
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
		FResult: res,
		FReturn: ret,
	})
}
