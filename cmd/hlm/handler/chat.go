package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hlm/database"
	"github.com/number571/go-peer/cmd/hlm/settings"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/settings"
	"github.com/number571/go-peer/cmd/hls/hlc"
	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

const (
	maxMessages = 1000
)

type sMessages struct {
	FMessages []sMessage
}

type sMessage struct {
	FIsIncoming bool
	FMessage    string
}

func FriendsChatPage(client hlc.IClient, db database.IKeyValueDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/friends/chat" {
			NotFoundPage(w, r)
			return
		}

		aliasName := r.URL.Query().Get("alias_name")
		if aliasName == "" {
			fmt.Fprint(w, "alias name is null")
			return
		}

		friends, err := client.GetFriends()
		if err != nil {
			fmt.Fprint(w, "error: read friends")
			return
		}

		pubKey, ok := friends[aliasName]
		if !ok {
			fmt.Fprint(w, "undefined public key by alias name")
			return
		}

		r.ParseForm()

		switch r.FormValue("method") {
		case http.MethodPost:
			msg := strings.TrimSpace(r.FormValue("input_message"))
			if msg == "" {
				fmt.Fprint(w, "error: message is null")
				return
			}

			res, err := client.Request(
				pubKey,
				hls_network.NewRequest("POST", settings.CTitlePattern, "/push").
					WithHead(map[string]string{
						"Content-Type": "application/json",
					}).
					WithBody([]byte(msg)),
			)
			if err != nil {
				fmt.Fprint(w, "error: push message to network")
				return
			}

			resp := &hls_settings.SResponse{}
			if err := json.Unmarshal(res, resp); err != nil {
				fmt.Fprint(w, "error: receive response")
				return
			}

			if resp.FResult != settings.CTitlePattern {
				fmt.Fprint(w, "error: invalid response")
				return
			}

			err = db.Push(pubKey, database.NewMessage(false, msg))
			if err != nil {
				fmt.Fprint(w, "error: add message to database")
				return
			}
			http.Redirect(w, r,
				fmt.Sprintf("/friends/chat?alias_name=%s", aliasName),
				http.StatusSeeOther)
			return
		}

		start := uint64(0)
		size := db.Size(pubKey)
		if size > maxMessages {
			start = size - maxMessages
		}

		msgs, err := db.Load(pubKey, start, size)
		if err != nil {
			fmt.Fprint(w, "error: read database")
			return
		}

		res := &sMessages{FMessages: make([]sMessage, 0, len(msgs))}
		for _, msg := range msgs {
			res.FMessages = append(res.FMessages, sMessage{
				FIsIncoming: msg.IsIncoming(),
				FMessage:    msg.GetMessage(),
			})
		}

		t, err := template.ParseFiles(
			hlm_settings.CPathTemplates+"index.html",
			hlm_settings.CPathTemplates+"chat.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		t.Execute(w, res)
	}
}
