package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
)

type sFriends struct {
	*sTemplate
	FFriends []string
}

func FriendsPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/friends" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		pR.ParseForm()

		client := getClient(pCfg)

		switch pR.FormValue("method") {
		case http.MethodPost:
			aliasName := strings.TrimSpace(pR.FormValue("alias_name"))
			pubStrKey := strings.TrimSpace(pR.FormValue("public_key"))
			if aliasName == hlm_settings.CIamAliasName || aliasName == "" || pubStrKey == "" {
				pLogger.PushWarn(logBuilder.WithMessage("get_alias_name"))
				fmt.Fprint(pW, "error: host or port is null")
				return
			}
			pubKey := asymmetric.LoadRSAPubKey(pubStrKey)
			if pubKey == nil {
				pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
				fmt.Fprint(pW, "error: public key is nil")
				return
			}
			if err := client.AddFriend(aliasName, pubKey); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("add_friend"))
				fmt.Fprint(pW, "error: add friend")
				return
			}
		case http.MethodDelete:
			aliasName := strings.TrimSpace(pR.FormValue("alias_name"))
			if aliasName == hlm_settings.CIamAliasName || aliasName == "" {
				pLogger.PushWarn(logBuilder.WithMessage("get_alias_name"))
				fmt.Fprint(pW, "error: alias_name is null")
				return
			}
			if err := client.DelFriend(aliasName); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("del_friend"))
				fmt.Fprint(pW, "error: del friend")
				return
			}
		}

		friends, err := client.GetFriends()
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
			fmt.Fprint(pW, "error: read friends")
			return
		}

		result := new(sFriends)
		result.sTemplate = getTemplate(pCfg)
		result.FFriends = make([]string, 0, len(friends)+1) // +1 CIamAliasName

		friendsList := make([]string, 0, len(friends))
		for aliasName := range friends {
			friendsList = append(friendsList, aliasName)
		}
		sort.Strings(friendsList)

		result.FFriends = append(result.FFriends, hlm_settings.CIamAliasName) // in top
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
		t.Execute(pW, result)
	}
}
