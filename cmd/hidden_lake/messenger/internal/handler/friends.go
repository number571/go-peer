package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type sFriends struct {
	*state.STemplateState
	FFriends []string
}

func FriendsPage(pStateManager state.IStateManager) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/friends" {
			NotFoundPage(pStateManager)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		pR.ParseForm()

		switch pR.FormValue("method") {
		case http.MethodPost:
			aliasName := strings.TrimSpace(pR.FormValue("alias_name"))
			pubStrKey := strings.TrimSpace(pR.FormValue("public_key"))
			if aliasName == "" || pubStrKey == "" {
				fmt.Fprint(pW, "error: host or port is null")
				return
			}
			pubKey := asymmetric.LoadRSAPubKey(pubStrKey)
			if pubKey == nil {
				fmt.Fprint(pW, "error: public key is nil")
				return
			}
			if err := pStateManager.AddFriend(aliasName, pubKey); err != nil {
				fmt.Fprint(pW, "error: add friend")
				return
			}
		case http.MethodDelete:
			aliasName := strings.TrimSpace(pR.FormValue("alias_name"))
			if aliasName == "" {
				fmt.Fprint(pW, "error: alias_name is null")
				return
			}
			if err := pStateManager.DelFriend(aliasName); err != nil {
				fmt.Fprint(pW, "error: del friend")
				return
			}
		}
		res, err := pStateManager.GetClient().Service().GetFriends()
		if err != nil {
			fmt.Fprint(pW, "error: read friends")
			return
		}

		result := new(sFriends)
		result.STemplateState = pStateManager.GetTemplate()
		result.FFriends = make([]string, 0, len(res))

		for aliasName := range res {
			result.FFriends = append(result.FFriends, aliasName)
		}
		sort.Strings(result.FFriends)

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"friends.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}
		t.Execute(pW, result)
	}
}
