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

func HandleConfigFriendsAPI(pWrapper config.IWrapper, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		var vFriend pkg_settings.SFriend

		switch pR.Method {
		case http.MethodGet:
			friends := pWrapper.GetConfig().GetFriends()
			listFriends := make([]string, 0, len(friends))
			for name, pubKey := range friends {
				listFriends = append(listFriends, fmt.Sprintf("%s:%s", name, pubKey.ToString()))
			}
			api.Response(pW, http.StatusOK, strings.Join(listFriends, ","))
			return
		case http.MethodPost, http.MethodDelete:
			// next
		default:
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vFriend); err != nil {
			api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		aliasName := strings.TrimSpace(vFriend.FAliasName)
		if aliasName == "" {
			api.Response(pW, http.StatusTeapot, "failed: load alias name")
			return
		}

		friends := pWrapper.GetConfig().GetFriends()

		switch pR.Method {
		case http.MethodPost:
			if _, ok := friends[aliasName]; ok {
				api.Response(pW, http.StatusNotAcceptable, "failed: friend already exist")
				return
			}

			pubKey := asymmetric.LoadRSAPubKey(vFriend.FPublicKey)
			if pubKey == nil {
				api.Response(pW, http.StatusBadRequest, "failed: load public key")
				return
			}

			friends[aliasName] = pubKey
			if err := pWrapper.GetEditor().UpdateFriends(friends); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: update friends")
				return
			}

			pNode.GetListPubKeys().AddPubKey(pubKey)
			api.Response(pW, http.StatusOK, "success: update friends")
		case http.MethodDelete:
			pubKey, ok := friends[aliasName]
			if !ok {
				api.Response(pW, http.StatusNotFound, "failed: friend does not exist")
				return
			}

			delete(friends, aliasName)
			if err := pWrapper.GetEditor().UpdateFriends(friends); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: delete friend")
				return
			}

			pNode.GetListPubKeys().DelPubKey(pubKey)
			api.Response(pW, http.StatusOK, "success: delete friend")
		}
	}
}
