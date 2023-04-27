package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleConfigFriendsAPI(pWrapper config.IWrapper, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		var vFriend pkg_settings.SFriend

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			api.Response(pW, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			friends := pWrapper.GetConfig().GetFriends()
			listFriends := make([]string, 0, len(friends))
			for name, pubKey := range friends {
				listFriends = append(listFriends, fmt.Sprintf("%s:%s", name, pubKey.ToString()))
			}
			api.Response(pW, pkg_settings.CErrorNone, strings.Join(listFriends, ","))
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vFriend); err != nil {
			api.Response(pW, pkg_settings.CErrorDecode, "failed: decode request")
			return
		}

		aliasName := strings.TrimSpace(vFriend.FAliasName)
		if aliasName == "" {
			api.Response(pW, pkg_settings.CErrorValue, "failed: load alias name")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			pubKey := asymmetric.LoadRSAPubKey(vFriend.FPublicKey)
			if pubKey == nil {
				api.Response(pW, pkg_settings.CErrorPubKey, "failed: load public key")
				return
			}

			friends := pWrapper.GetConfig().GetFriends()
			if _, ok := friends[aliasName]; ok {
				api.Response(pW, pkg_settings.CErrorExist, "failed: friend already exist")
				return
			}

			friends[aliasName] = pubKey
			if err := pWrapper.GetEditor().UpdateFriends(friends); err != nil {
				api.Response(pW, pkg_settings.CErrorAction, "failed: update friends")
				return
			}

			pNode.GetListPubKeys().AddPubKey(pubKey)
			api.Response(pW, pkg_settings.CErrorNone, "success: update friends")
		case http.MethodDelete:
			friends := pWrapper.GetConfig().GetFriends()
			pubKey, ok := friends[aliasName]
			if !ok {
				api.Response(pW, pkg_settings.CErrorNotExist, "failed: friend does not exist")
				return
			}

			delete(friends, aliasName)
			if err := pWrapper.GetEditor().UpdateFriends(friends); err != nil {
				api.Response(pW, pkg_settings.CErrorAction, "failed: delete friend"+err.Error())
				return
			}

			pNode.GetListPubKeys().DelPubKey(pubKey)
			api.Response(pW, pkg_settings.CErrorNone, "success: delete friend")
		}
	}
}
