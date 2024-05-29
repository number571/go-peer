package handler

import (
	"context"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/web"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
)

type sFriends struct {
	*sTemplate
	FFriends []string
}

func FriendsPage(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.URL.Path != "/friends" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		if err := pR.ParseForm(); err != nil {
			ErrorPage(pLogger, pCfg, "parse_form", "parse form")(pW, pR)
			return
		}

		switch pR.FormValue("method") {
		case http.MethodPost:
			pubStrKey := strings.TrimSpace(pR.FormValue("public_key"))
			aliasName := strings.TrimSpace(pR.FormValue("alias_name")) // may be nil

			if pubStrKey == "" {
				ErrorPage(pLogger, pCfg, "public_key_nil", "public key is nil")(pW, pR)
				return
			}

			pubKey := asymmetric.LoadRSAPubKey(pubStrKey)
			if pubKey == nil {
				ErrorPage(pLogger, pCfg, "decode_public_key", "failed decode public key")(pW, pR)
				return
			}

			if aliasName == "" {
				// get hash of public key as alias_name
				aliasName = pubKey.GetHasher().ToString()
			}

			if err := pHlsClient.AddFriend(pCtx, aliasName, pubKey); err != nil {
				ErrorPage(pLogger, pCfg, "add_friend", "add friend")(pW, pR)
				return
			}
		case http.MethodDelete:
			aliasName := strings.TrimSpace(pR.FormValue("alias_name"))
			if aliasName == "" {
				ErrorPage(pLogger, pCfg, "get_alias_name", "alias_name is nil")(pW, pR)
				return
			}

			if err := pHlsClient.DelFriend(pCtx, aliasName); err != nil {
				ErrorPage(pLogger, pCfg, "del_friend", "delete friend")(pW, pR)
				return
			}
		}

		friends, err := pHlsClient.GetFriends(pCtx)
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_friends", "read friends")(pW, pR)
			return
		}

		result := new(sFriends)
		result.sTemplate = getTemplate(pCfg)
		result.FFriends = make([]string, 0, len(friends))

		friendsList := make([]string, 0, len(friends))
		for aliasName := range friends {
			friendsList = append(friendsList, aliasName)
		}
		sort.Strings(friendsList)

		result.FFriends = append(result.FFriends, friendsList...)

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"friends.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = t.Execute(pW, result)
	}
}
