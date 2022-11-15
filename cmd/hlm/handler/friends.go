package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	hlm_settings "github.com/number571/go-peer/cmd/hlm/settings"
	"github.com/number571/go-peer/cmd/hls/hlc"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

type sFriends struct {
	FFriends []string
}

func FriendsPage(client hlc.IClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/friends" {
			NotFoundPage(w, r)
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
			err := client.AddFriend(aliasName, pubKey)
			if err != nil {
				fmt.Fprint(w, "error: add connection")
				return
			}
		case http.MethodDelete:
			aliasName := strings.TrimSpace(r.FormValue("alias_name"))
			if aliasName == "" {
				fmt.Fprint(w, "error: alias_name is null")
				return
			}
			err := client.DelFriend(aliasName)
			if err != nil {
				fmt.Fprint(w, "error: del connection")
				return
			}
		}

		res, err := client.GetFriends()
		if err != nil {
			fmt.Fprint(w, "error: read friends")
			return
		}

		result := new(sFriends)
		result.FFriends = make([]string, 0, len(res))
		for aliasName := range res {
			result.FFriends = append(result.FFriends, aliasName)
		}
		sort.Strings(result.FFriends)

		t, err := template.ParseFiles(
			hlm_settings.CPathTemplates+"index.html",
			hlm_settings.CPathTemplates+"friends.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}
		t.Execute(w, result)
	}
}
