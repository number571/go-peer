package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
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
	*state.STemplateState
	FAddress  sChatAddress
	FMessages []sChatMessage
}

func FriendsChatPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/friends/chat" {
			NotFoundPage(s)(w, r)
			return
		}

		if !s.IsActive() {
			http.Redirect(w, r, "/sign/in", http.StatusFound)
			return
		}

		aliasName := r.URL.Query().Get("alias_name")
		if aliasName == "" {
			fmt.Fprint(w, "alias name is null")
			return
		}

		var (
			client = s.GetClient().Service()
			db     = s.GetWrapperDB().Get()
		)

		myPubKey, err := client.GetPubKey()
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

			res, err := client.DoRequest(
				friendPubKey,
				request.NewRequest(http.MethodPost, hlm_settings.CTitlePattern, "/push").
					WithHead(map[string]string{
						"Content-Type": "application/json",
					}).
					WithBody([]byte(msg)),
			)
			if err != nil {
				fmt.Fprint(w, "error: push message to network")
				return
			}

			resp := &api.SResponse{}
			if err := json.Unmarshal(res, resp); err != nil {
				fmt.Fprint(w, "error: receive response")
				return
			}

			if resp.FResult != hlm_settings.CTitlePattern {
				fmt.Fprint(w, "error: invalid response")
				return
			}

			uid := random.NewStdPRNG().Bytes(hashing.CSHA256Size)
			err = db.Push(rel, database.NewMessage(false, msg, uid))
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

		clientPubKey, err := client.GetPubKey()
		if err != nil {
			fmt.Fprint(w, "error: read public key")
			return
		}

		res := &sChatMessages{
			STemplateState: s.GetTemplate(),
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

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"chat.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		t.Execute(w, res)
	}
}
