package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hlm/internal/database"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/internal/settings"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
	"github.com/number571/go-peer/cmd/hls/pkg/request"
	hls_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
)

const (
	cChatLimitMessages = 1000
)

type sChatMessage struct {
	FIsIncoming bool
	FTimestamp  string
	FMessage    string
}
type sChatAddress struct {
	FClient string
	FFriend string
}
type sChatMessages struct {
	*sTemplateData
	FAddress  sChatAddress
	FMessages []sChatMessage
}

func FriendsChatPage(wDB database.IWrapperDB, client hls_client.IClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := wDB.Get()
		if r.URL.Path != "/friends/chat" {
			NotFoundPage(db)(w, r)
			return
		}

		if db == nil {
			http.Redirect(w, r, "/sign/in", http.StatusFound)
			return
		}

		aliasName := r.URL.Query().Get("alias_name")
		if aliasName == "" {
			fmt.Fprint(w, "alias name is null")
			return
		}

		myPubKey, err := client.PubKey()
		if err != nil {
			fmt.Fprint(w, "error: read public key")
			return
		}

		friends, err := client.GetFriends()
		if err != nil {
			fmt.Fprint(w, "error: read friends")
			return
		}

		friendPubKey, ok := friends[aliasName]
		if !ok {
			fmt.Fprint(w, "undefined public key by alias name")
			return
		}

		rel := database.NewRelation(myPubKey, friendPubKey)
		r.ParseForm()

		switch r.FormValue("method") {
		case http.MethodPost:
			msg := strings.TrimSpace(r.FormValue("input_message"))
			if msg == "" {
				fmt.Fprint(w, "error: message is null")
				return
			}

			res, err := client.Request(
				friendPubKey,
				request.NewRequest("POST", hlm_settings.CTitlePattern, "/push").
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

			if resp.FResult != hlm_settings.CTitlePattern {
				fmt.Fprint(w, "error: invalid response")
				return
			}

			err = db.Push(rel, database.NewMessage(false, msg))
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
		size := db.Size(rel)
		if size > cChatLimitMessages {
			start = size - cChatLimitMessages
		}

		msgs, err := db.Load(rel, start, size)
		if err != nil {
			fmt.Fprint(w, "error: read database")
			return
		}

		clientPubKey, err := client.PubKey()
		if err != nil {
			fmt.Fprint(w, "error: read public key")
			return
		}

		res := &sChatMessages{
			sTemplateData: newTemplateData(db),
			FAddress: sChatAddress{
				FClient: clientPubKey.Address().String(),
				FFriend: friendPubKey.Address().String(),
			},
			FMessages: make([]sChatMessage, 0, len(msgs)),
		}
		for _, msg := range msgs {
			res.FMessages = append(res.FMessages, sChatMessage{
				FIsIncoming: msg.IsIncoming(),
				FTimestamp:  msg.GetTimestamp(),
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
