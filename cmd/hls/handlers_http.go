package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"

	"github.com/number571/go-peer/cmd/hls/config"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	payload_adapter "github.com/number571/go-peer/modules/network/anonymity/adapters/payload"
)

func handleIndexHTTP(w http.ResponseWriter, r *http.Request) {
	response(w, hls_settings.CErrorNone, hls_settings.CTitlePattern)
}

func handlePubKeyHTTP(w http.ResponseWriter, r *http.Request) {
	response(w, hls_settings.CErrorNone, gNode.Queue().Client().PubKey().String())
}

func handleOnlineHTTP(w http.ResponseWriter, r *http.Request) {
	var vConnect hls_settings.SConnect

	if r.Method != "GET" && r.Method != "DELETE" {
		response(w, hls_settings.CErrorMethod, "failed: incorrect method")
		return
	}

	if r.Method == "GET" {
		conns := gNode.Network().Connections()
		inOnline := make([]string, 0, len(conns))
		for _, conn := range conns {
			inOnline = append(inOnline, conn.Socket().RemoteAddr().String())
		}
		sort.SliceStable(inOnline, func(i, j int) bool {
			return inOnline[i] < inOnline[j]
		})
		response(w, hls_settings.CErrorNone, strings.Join(inOnline, ","))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&vConnect); err != nil {
		response(w, hls_settings.CErrorDecode, "failed: decode request")
		return
	}

	if err := gNode.Network().Disconnect(vConnect.FConnect); err != nil {
		response(w, hls_settings.CErrorNone, "failed: delete online connection")
		return
	}
	response(w, hls_settings.CErrorNone, "success: delete online connection")
}

func handleConnectsHTTP(w http.ResponseWriter, r *http.Request) {
	var vConnect hls_settings.SConnect

	if r.Method != "GET" && r.Method != "PUT" && r.Method != "DELETE" {
		response(w, hls_settings.CErrorMethod, "failed: incorrect method")
		return
	}

	if r.Method == "GET" {
		response(w, hls_settings.CErrorNone, strings.Join(gConfig.Connections(), ","))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&vConnect); err != nil {
		response(w, hls_settings.CErrorDecode, "failed: decode request")
		return
	}

	switch r.Method {
	case "PUT":
		connects := append(gConfig.Connections(), vConnect.FConnect)
		if err := gEditor.UpdateConnections(connects); err != nil {
			response(w, hls_settings.CErrorAction, "failed: update connections")
			return
		}
		gNode.Network().Connect(vConnect.FConnect)
		response(w, hls_settings.CErrorNone, "success: update connections")
	case "DELETE":
		connects := deleteConnect(gConfig, vConnect.FConnect)
		if err := gEditor.UpdateConnections(connects); err != nil {
			response(w, hls_settings.CErrorAction, "failed: delete connection")
			return
		}
		gNode.Network().Disconnect(vConnect.FConnect)
		response(w, hls_settings.CErrorNone, "success: delete connection")
	}
}

func handleFriendsHTTP(w http.ResponseWriter, r *http.Request) {
	var vFriend hls_settings.SFriend

	if r.Method != "GET" && r.Method != "PUT" && r.Method != "DELETE" {
		response(w, hls_settings.CErrorMethod, "failed: incorrect method")
		return
	}

	if r.Method == "GET" {
		listPubKeys := gNode.F2F().List()
		listPubKeysStr := make([]string, 0, len(listPubKeys))
		for _, pubKey := range listPubKeys {
			listPubKeysStr = append(listPubKeysStr, pubKey.String())
		}
		response(w, hls_settings.CErrorNone, strings.Join(listPubKeysStr, ","))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&vFriend); err != nil {
		response(w, hls_settings.CErrorDecode, "failed: decode request")
		return
	}

	pubKey := asymmetric.LoadRSAPubKey(vFriend.FFriend)
	if pubKey == nil {
		response(w, hls_settings.CErrorPubKey, "failed: load public key")
		return
	}

	switch r.Method {
	case "PUT":
		friends := append(gConfig.Friends(), pubKey)
		if err := gEditor.UpdateFriends(friends); err != nil {
			response(w, hls_settings.CErrorAction, "failed: update friends")
			return
		}
		gNode.F2F().Append(pubKey)
		response(w, hls_settings.CErrorNone, "success: update friends")
	case "DELETE":
		friends := deleteFriend(gConfig, pubKey)
		if err := gEditor.UpdateFriends(friends); err != nil {
			response(w, hls_settings.CErrorAction, "failed: delete friend"+err.Error())
			return
		}
		gNode.F2F().Remove(pubKey)
		response(w, hls_settings.CErrorNone, "success: delete friend")
	}
}

func handlePushHTTP(w http.ResponseWriter, r *http.Request) {
	var vPush hls_settings.SPush

	if r.Method != "POST" && r.Method != "PUT" {
		response(w, hls_settings.CErrorMethod, "failed: incorrect method")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&vPush); err != nil {
		response(w, hls_settings.CErrorDecode, "failed: decode request")
		return
	}

	pubKey := asymmetric.LoadRSAPubKey(vPush.FReceiver)
	if pubKey == nil {
		response(w, hls_settings.CErrorPubKey, "failed: load public key")
		return
	}

	data := encoding.HexDecode(vPush.FHexData)
	if data == nil {
		response(w, hls_settings.CErrorPubKey, "failed: decode hex format data")
		return
	}

	switch r.Method {
	case "PUT":
		msg, err := gNode.Queue().Client().Encrypt(
			pubKey,
			payload_adapter.NewPayload(hls_settings.CHeaderHLS, data),
		)
		if err != nil {
			response(w, hls_settings.CErrorMessage, "failed: encrypt message with data")
			return
		}
		if err := gNode.Broadcast(msg); err != nil {
			response(w, hls_settings.CErrorBroadcast, "failed: broadcast message")
			return
		}
		response(w, hls_settings.CErrorNone, "success: broadcast")
		return
	case "POST":
		resp, err := gNode.Request(
			pubKey,
			payload_adapter.NewPayload(hls_settings.CHeaderHLS, data),
		)
		if err != nil {
			response(w, hls_settings.CErrorResponse, "failed: response message")
			return
		}
		response(w, hls_settings.CErrorNone, encoding.HexEncode(resp))
		return
	}
}

func response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", hls_settings.CContentType)
	json.NewEncoder(w).Encode(&hls_settings.SResponse{
		FResult: res,
		FReturn: ret,
	})
}

func deleteConnect(cfg config.IConfig, connect string) []string {
	connects := gConfig.Connections()
	result := make([]string, 0, len(connects))
	for _, conn := range connects {
		if conn == connect {
			continue
		}
		result = append(result, conn)
	}
	return result
}

func deleteFriend(cfg config.IConfig, friend asymmetric.IPubKey) []asymmetric.IPubKey {
	friends := gConfig.Friends()
	result := make([]asymmetric.IPubKey, 0, len(friends))
	for _, f := range friends {
		if f.Address().String() == friend.Address().String() {
			continue
		}
		result = append(result, f)
	}
	return result
}
