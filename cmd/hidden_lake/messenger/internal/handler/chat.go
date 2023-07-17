package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/errors"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
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

		recvPubKey, err := getReceiverPubKey(client, myPubKey, aliasName)
		if err != nil {
			fmt.Fprint(pW, err.Error())
			return
		}

		rel := database.NewRelation(myPubKey, recvPubKey)
		pR.ParseForm()

		switch pR.FormValue("method") {
		case http.MethodPost:
			rawMsg := strings.TrimSpace(pR.FormValue("input_message"))

			// ReplaceTextToEmoji not used in trySendMessage
			// because can be exists emoji (🕵️‍♂️ => 🕵️♂️) more than 4 bytes (rune)
			// that creates more than 1 emoji after use GetOnlyWritableCharacters
			msg := utils.GetOnlyWritableCharacters(rawMsg)
			if msg == "" {
				fmt.Fprint(pW, "error: message is null")
				return
			}

			if err := trySendMessage(client, recvPubKey, myPubKey, msg); err != nil {
				fmt.Fprint(pW, "error: push message to network")
				return
			}

			uid := random.NewStdPRNG().GetBytes(hashing.CSHA256Size)
			dbMsg := database.NewMessage(false, utils.ReplaceTextToEmoji(msg), uid)

			if err := db.Push(rel, dbMsg); err != nil {
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
		if size > hlm_settings.CChatLimitMessages {
			start = size - hlm_settings.CChatLimitMessages
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
				FFriend: recvPubKey.GetAddress().ToString(),
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

func trySendMessage(client client.IClient, recvPubKey, myPubKey asymmetric.IPubKey, msg string) error {
	// if the sender = receiver then there is no need to send a message to the network
	if myPubKey.GetAddress().ToString() == recvPubKey.GetAddress().ToString() {
		return nil
	}

	return client.BroadcastRequest(
		recvPubKey,
		request.NewRequest(http.MethodPost, pkg_settings.CTitlePattern, hlm_settings.CPushPath).
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(msg)),
	)
}

func getReceiverPubKey(client client.IClient, myPubKey asymmetric.IPubKey, aliasName string) (asymmetric.IPubKey, error) {
	friends, err := client.GetFriends()
	if err != nil {
		return nil, errors.NewError("error: read friends")
	}

	if aliasName == hlm_settings.CIamAliasName {
		return myPubKey, nil
	}

	friendPubKey, ok := friends[aliasName]
	if !ok {
		return nil, errors.NewError("undefined public key by alias name")
	}

	return friendPubKey, nil
}
