package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
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

func FriendsChatPage(pStateManager state.IStateManager) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/friends/chat" {
			NotFoundPage(pStateManager)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		aliasName := pR.URL.Query().Get("alias_name")
		if aliasName == "" {
			fmt.Fprint(pW, "alias name is null")
			return
		}

		db := pStateManager.GetWrapperDB().Get()
		if db == nil {
			api.Response(pW, http.StatusForbidden, "failed: database closed")
			return
		}

		client := pStateManager.GetClient().Service()
		myPubKey, err := client.GetPubKey()
		if err != nil {
			fmt.Fprint(pW, "error: read public key")
			return
		}

		friends, err := client.GetFriends()
		if err != nil {
			fmt.Fprint(pW, "error: read friends")
			return
		}

		friendPubKey, ok := friends[aliasName]
		if !ok {
			fmt.Fprint(pW, "undefined public key by alias name")
			return
		}

		rel := database.NewRelation(myPubKey, friendPubKey)
		pR.ParseForm()

		switch pR.FormValue("method") {
		case http.MethodPost:
			msg := strings.TrimSpace(pR.FormValue("input_message"))
			if msg == "" {
				fmt.Fprint(pW, "error: message is null")
				return
			}

			resp, err := client.FetchRequest(
				friendPubKey,
				request.NewRequest(http.MethodPost, pkg_settings.CTitlePattern, hlm_settings.CPushPath).
					WithHead(map[string]string{
						"Content-Type": "application/json",
					}).
					WithBody([]byte(msg)),
			)
			if err != nil {
				fmt.Fprint(pW, "error: push message to network")
				return
			}

			if string(resp.GetBody()) != pkg_settings.CTitlePattern {
				fmt.Fprint(pW, "error: invalid response")
				return
			}

			uid := random.NewStdPRNG().GetBytes(hashing.CSHA256Size)
			err = db.Push(rel, database.NewMessage(false, msg, uid))
			if err != nil {
				fmt.Fprint(pW, "error: add message to database")
				return
			}

			http.Redirect(pW, pR,
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
			fmt.Fprint(pW, "error: read database")
			return
		}

		clientPubKey, err := client.GetPubKey()
		if err != nil {
			fmt.Fprint(pW, "error: read public key")
			return
		}

		res := &sChatMessages{
			STemplateState: pStateManager.GetTemplate(),
			FAddress: sChatAddress{
				FClient: clientPubKey.GetAddress().ToString(),
				FFriend: friendPubKey.GetAddress().ToString(),
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

		t.Execute(pW, res)
	}
}
