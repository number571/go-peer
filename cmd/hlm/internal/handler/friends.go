package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/number571/go-peer/cmd/hlm/internal/app/state"
	"github.com/number571/go-peer/cmd/hlm/web"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type sFriends struct {
	*state.STemplateState
	FFriends []string
}

func FriendsPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/friends" {
			NotFoundPage(s)(w, r)
			return
		}

		if !s.IsActive() {
			http.Redirect(w, r, "/sign/in", http.StatusFound)
			return
		}

		r.ParseForm()

		switch r.FormValue("method") {
		case http.MethodPost:
			aliasName := strings.TrimSpace(r.FormValue("alias_name"))
			pubStrKey := strings.TrimSpace(r.FormValue("public_key"))
			if aliasName == "" || pubStrKey == "" {
				fmt.Fprint(w, "error: host or port is null")
				return
			}
			pubKey := asymmetric.LoadRSAPubKey(pubStrKey)
			if pubKey == nil {
				fmt.Fprint(w, "error: public key is nil")
				return
			}
			if err := s.AddFriend(aliasName, pubKey); err != nil {
				fmt.Fprint(w, "error: add friend")
				return
			}
		case http.MethodDelete:
			aliasName := strings.TrimSpace(r.FormValue("alias_name"))
			if aliasName == "" {
				fmt.Fprint(w, "error: alias_name is null")
				return
			}
			if err := s.DelFriend(aliasName); err != nil {
				fmt.Fprint(w, "error: del friend")
				return
			}
		}
		res, err := s.GetClient().GetFriends()
		if err != nil {
			fmt.Fprint(w, "error: read friends")
			return
		}

		result := new(sFriends)
		result.STemplateState = s.GetTemplate()
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
		t.Execute(w, result)
	}
}
