package handler

import (
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/web"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
)

type sFriends struct {
	*sTemplate
	FFriends []string
}

func FriendsPage(pLogger logger.ILogger, pWrapper config.IWrapper) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		cfg := pWrapper.GetConfig()
		cfgEditor := pWrapper.GetEditor()

		if pR.URL.Path != "/friends" {
			NotFoundPage(pLogger, cfg)(pW, pR)
			return
		}

		pR.ParseForm()

		client := getClient(cfg)
		secretKeys := cfg.GetSecretKeys()

		switch pR.FormValue("method") {
		case http.MethodPost:
			pubStrKey := strings.TrimSpace(pR.FormValue("public_key"))
			aliasName := strings.TrimSpace(pR.FormValue("alias_name")) // may be nil
			secretKey := strings.TrimSpace(pR.FormValue("secret_key")) // may be nil

			if pubStrKey == "" {
				ErrorPage(pLogger, cfg, "public_key_nil", "public key is nil")(pW, pR)
				return
			}

			pubKey := asymmetric.LoadRSAPubKey(pubStrKey)
			if pubKey == nil {
				ErrorPage(pLogger, cfg, "decode_public_key", "failed decode public key")(pW, pR)
				return
			}

			if aliasName == "" {
				// get hash of public key as alias_name
				aliasName = pubKey.GetHasher().ToString()
			}

			if err := client.AddFriend(aliasName, pubKey); err != nil {
				ErrorPage(pLogger, cfg, "add_friend", "add friend")(pW, pR)
				return
			}

			secretKeys[aliasName] = secretKey
			if err := cfgEditor.UpdateSecretKeys(secretKeys); err != nil {
				ErrorPage(pLogger, cfg, "add_secret_key", "add secret key")(pW, pR)
				return
			}
		case http.MethodDelete:
			aliasName := strings.TrimSpace(pR.FormValue("alias_name"))
			if aliasName == "" {
				ErrorPage(pLogger, cfg, "get_alias_name", "alias_name is nil")(pW, pR)
				return
			}

			if err := client.DelFriend(aliasName); err != nil {
				ErrorPage(pLogger, cfg, "del_friend", "delete friend")(pW, pR)
				return
			}

			delete(secretKeys, aliasName)
			if err := cfgEditor.UpdateSecretKeys(secretKeys); err != nil {
				ErrorPage(pLogger, cfg, "del_secret_key", "delete secret key")(pW, pR)
				return
			}
		}

		friends, err := client.GetFriends()
		if err != nil {
			ErrorPage(pLogger, cfg, "get_friends", "read friends")(pW, pR)
			return
		}

		result := new(sFriends)
		result.sTemplate = getTemplate(cfg)
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
		t.Execute(pW, result)
	}
}
