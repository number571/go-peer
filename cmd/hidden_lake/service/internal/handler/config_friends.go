package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleConfigFriendsAPI(wrapper config.IWrapper, node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vFriend pkg_settings.SFriend

		if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodDelete {
			api.Response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if r.Method == http.MethodGet {
			friends := wrapper.Config().Friends()
			listFriends := make([]string, 0, len(friends))
			for name, pubKey := range friends {
				listFriends = append(listFriends, fmt.Sprintf("%s:%s", name, pubKey.ToString()))
			}
			api.Response(w, pkg_settings.CErrorNone, strings.Join(listFriends, ","))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&vFriend); err != nil {
			api.Response(w, pkg_settings.CErrorDecode, "failed: decode request")
			return
		}

		aliasName := strings.TrimSpace(vFriend.FAliasName)
		if aliasName == "" {
			api.Response(w, pkg_settings.CErrorValue, "failed: load alias name")
			return
		}

		switch r.Method {
		case http.MethodPost:
			pubKey := asymmetric.LoadRSAPubKey(vFriend.FPublicKey)
			if pubKey == nil {
				api.Response(w, pkg_settings.CErrorPubKey, "failed: load public key")
				return
			}

			friends := wrapper.Config().Friends()
			if _, ok := friends[aliasName]; ok {
				api.Response(w, pkg_settings.CErrorExist, "failed: friend already exist")
				return
			}

			friends[aliasName] = pubKey
			if err := wrapper.Editor().UpdateFriends(friends); err != nil {
				api.Response(w, pkg_settings.CErrorAction, "failed: update friends")
				return
			}

			node.GetListPubKeys().AddPubKey(pubKey)
			api.Response(w, pkg_settings.CErrorNone, "success: update friends")
		case http.MethodDelete:
			friends := wrapper.Config().Friends()
			pubKey, ok := friends[aliasName]
			if !ok {
				api.Response(w, pkg_settings.CErrorNotExist, "failed: friend does not exist")
				return
			}

			delete(friends, aliasName)
			if err := wrapper.Editor().UpdateFriends(friends); err != nil {
				api.Response(w, pkg_settings.CErrorAction, "failed: delete friend"+err.Error())
				return
			}

			node.GetListPubKeys().DelPubKey(pubKey)
			api.Response(w, pkg_settings.CErrorNone, "success: delete friend")
		}
	}
}
