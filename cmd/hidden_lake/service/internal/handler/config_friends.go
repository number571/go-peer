package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleConfigFriendsAPI(pWrapper config.IWrapper, pLogger logger.ILogger, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		var vFriend pkg_settings.SFriend

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			friends := pWrapper.GetConfig().GetFriends()

			listFriends := make([]pkg_settings.SFriend, 0, len(friends))
			for name, pubKey := range friends {
				listFriends = append(listFriends, pkg_settings.SFriend{
					FAliasName: name,
					FPublicKey: pubKey.ToString(),
				})
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, listFriends)
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vFriend); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		aliasName := strings.TrimSpace(vFriend.FAliasName)
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("get_alias_name"))
			api.Response(pW, http.StatusTeapot, "failed: load alias name")
			return
		}

		friends := pWrapper.GetConfig().GetFriends()

		switch pR.Method {
		case http.MethodPost:
			if _, ok := friends[aliasName]; ok {
				pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
				api.Response(pW, http.StatusNotAcceptable, "failed: friend already exist")
				return
			}

			pubKey := asymmetric.LoadRSAPubKey(vFriend.FPublicKey)
			if pubKey == nil {
				pLogger.PushWarn(logBuilder.WithMessage("decode_key"))
				api.Response(pW, http.StatusBadRequest, "failed: load public key")
				return
			}

			friends[aliasName] = pubKey
			if err := pWrapper.GetEditor().UpdateFriends(friends); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_friends"))
				api.Response(pW, http.StatusInternalServerError, "failed: update friends")
				return
			}

			pNode.GetListPubKeys().AddPubKey(pubKey)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: update friends")
			return

		case http.MethodDelete:
			pubKey, ok := friends[aliasName]
			if !ok {
				pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
				api.Response(pW, http.StatusNotFound, "failed: friend does not exist")
				return
			}

			delete(friends, aliasName)

			if err := pWrapper.GetEditor().UpdateFriends(friends); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_friends"))
				api.Response(pW, http.StatusInternalServerError, "failed: delete friend")
				return
			}

			pNode.GetListPubKeys().DelPubKey(pubKey)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: delete friend")
			return
		}
	}
}
