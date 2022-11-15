package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"

	"github.com/number571/go-peer/cmd/hls/config"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/network/anonymity"
	payload_adapter "github.com/number571/go-peer/modules/network/anonymity/adapters/payload"
)

func handleIndexHTTP(w http.ResponseWriter, r *http.Request) {
	response(w, hls_settings.CErrorNone, hls_settings.CTitlePattern)
}

func handlePubKeyHTTP(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response(w, hls_settings.CErrorNone, node.Queue().Client().PubKey().String())
	}
}

func handleConnectsHTTP(wrapper config.IWrapper, node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vConnect hls_settings.SConnect

		if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodDelete {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if r.Method == http.MethodGet {
			response(w, hls_settings.CErrorNone, strings.Join(wrapper.Config().Connections(), ","))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&vConnect); err != nil {
			response(w, hls_settings.CErrorDecode, "failed: decode request")
			return
		}

		switch r.Method {
		case http.MethodPost:
			connects := append(wrapper.Config().Connections(), vConnect.FConnect)
			if err := wrapper.Editor().UpdateConnections(connects); err != nil {
				response(w, hls_settings.CErrorAction, "failed: update connections")
				return
			}
			node.Network().Connect(vConnect.FConnect)
			response(w, hls_settings.CErrorNone, "success: update connections")
		case http.MethodDelete:
			connects := deleteConnect(wrapper.Config(), vConnect.FConnect)
			if err := wrapper.Editor().UpdateConnections(connects); err != nil {
				response(w, hls_settings.CErrorAction, "failed: delete connection")
				return
			}
			node.Network().Disconnect(vConnect.FConnect)
			response(w, hls_settings.CErrorNone, "success: delete connection")
		}
	}
}

func handleFriendsHTTP(wrapper config.IWrapper, node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vFriend hls_settings.SFriend

		if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodDelete {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if r.Method == http.MethodGet {
			friends := wrapper.Config().Friends()
			listFriends := make([]string, 0, len(friends))
			for name, pubKey := range friends {
				listFriends = append(listFriends, fmt.Sprintf("%s:%s", name, pubKey.String()))
			}
			response(w, hls_settings.CErrorNone, strings.Join(listFriends, ","))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&vFriend); err != nil {
			response(w, hls_settings.CErrorDecode, "failed: decode request")
			return
		}

		aliasName := strings.TrimSpace(vFriend.FAliasName)
		if aliasName == "" {
			response(w, hls_settings.CErrorValue, "failed: load alias name")
			return
		}

		switch r.Method {
		case http.MethodPost:
			pubKey := asymmetric.LoadRSAPubKey(vFriend.FPublicKey)
			if pubKey == nil {
				response(w, hls_settings.CErrorPubKey, "failed: load public key")
				return
			}

			friends := wrapper.Config().Friends()
			if _, ok := friends[aliasName]; ok {
				response(w, hls_settings.CErrorExist, "failed: friend already exist")
				return
			}

			friends[aliasName] = pubKey
			if err := wrapper.Editor().UpdateFriends(friends); err != nil {
				response(w, hls_settings.CErrorAction, "failed: update friends")
				return
			}

			node.F2F().Append(pubKey)
			response(w, hls_settings.CErrorNone, "success: update friends")
		case http.MethodDelete:
			friends := wrapper.Config().Friends()
			pubKey, ok := friends[aliasName]
			if !ok {
				response(w, hls_settings.CErrorNotExist, "failed: friend does not exist")
				return
			}

			delete(friends, aliasName)
			if err := wrapper.Editor().UpdateFriends(friends); err != nil {
				response(w, hls_settings.CErrorAction, "failed: delete friend"+err.Error())
				return
			}

			node.F2F().Remove(pubKey)
			response(w, hls_settings.CErrorNone, "success: delete friend")
		}
	}
}

func handleOnlineHTTP(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vConnect hls_settings.SConnect

		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if r.Method == http.MethodGet {
			conns := node.Network().Connections()
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

		if err := node.Network().Disconnect(vConnect.FConnect); err != nil {
			response(w, hls_settings.CErrorNone, "failed: delete online connection")
			return
		}
		response(w, hls_settings.CErrorNone, "success: delete online connection")
	}
}

func handlePushHTTP(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vPush hls_settings.SPush

		if r.Method != http.MethodPost && r.Method != http.MethodPut {
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
		case http.MethodPut:
			msg, err := node.Queue().Client().Encrypt(
				pubKey,
				payload_adapter.NewPayload(hls_settings.CHeaderHLS, data),
			)
			if err != nil {
				response(w, hls_settings.CErrorMessage, "failed: encrypt message with data")
				return
			}
			if err := node.Broadcast(msg); err != nil {
				response(w, hls_settings.CErrorBroadcast, "failed: broadcast message")
				return
			}
			response(w, hls_settings.CErrorNone, "success: broadcast")
			return
		case http.MethodPost:
			resp, err := node.Request(
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
}

func response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", hls_settings.CContentType)
	json.NewEncoder(w).Encode(&hls_settings.SResponse{
		FResult: res,
		FReturn: ret,
	})
}

func deleteConnect(cfg config.IConfig, connect string) []string {
	connects := cfg.Connections()
	result := make([]string, 0, len(connects))
	for _, conn := range connects {
		if conn == connect {
			continue
		}
		result = append(result, conn)
	}
	return result
}
