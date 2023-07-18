package handler

import (
	"fmt"
	"html/template"
	"io"
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
	pkg_client "github.com/number571/go-peer/pkg/client"
)

type sChatMessage struct {
	FIsIncoming  bool
	FMessageInfo utils.SMessageInfo
}
type sChatAddress struct {
	FAliasName string
	FFriend    string
}
type sChatMessages struct {
	*state.STemplateState
	FAddress  sChatAddress
	FMessages []sChatMessage
}

func FriendsChatPage(pStateManager state.IStateManager) http.HandlerFunc {
	msgSize := pStateManager.GetConfig().GetMessageSizeBytes()
	keySize := pStateManager.GetConfig().GetKeySizeBits()
	msgLimit := pkg_client.GetMessageLimit(msgSize, keySize)

	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/friends/chat" {
			NotFoundPage(pStateManager)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		pR.ParseForm()
		pR.ParseMultipartForm(10 << 20) // default value

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
		switch pR.FormValue("method") {
		case http.MethodPost, http.MethodPut:
			msgBytes, err := getMessageBytesSend(pStateManager, pR)
			if err != nil {
				fmt.Fprint(pW, errors.WrapError(err, "error: get message bytes"))
				return
			}

			if err := trySendMessage(client, recvPubKey, myPubKey, msgBytes, msgLimit); err != nil {
				fmt.Fprint(pW, errors.WrapError(err, "error: push message to network"))
				return
			}

			uid := random.NewStdPRNG().GetBytes(hashing.CSHA256Size)
			dbMsg := database.NewMessage(false, msgBytes, uid)

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

		messagesCap := pStateManager.GetConfig().GetMessagesCapacity()
		if size > messagesCap {
			start = size - messagesCap
		}

		msgs, err := db.Load(rel, start, size)
		if err != nil {
			fmt.Fprint(pW, "error: read database")
			return
		}

		res := &sChatMessages{
			STemplateState: pStateManager.GetTemplate(),
			FAddress: sChatAddress{
				FAliasName: aliasName,
				FFriend:    recvPubKey.GetAddress().ToString(),
			},
			FMessages: make([]sChatMessage, 0, len(msgs)),
		}

		for _, msg := range msgs {
			res.FMessages = append(res.FMessages, sChatMessage{
				FIsIncoming:  msg.IsIncoming(),
				FMessageInfo: getMessageInfo(msg.GetMessage(), msg.GetTimestamp()),
			})
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"chat.html",
		)
		if err != nil {
			panic(err)
		}

		t.Execute(pW, res)
	}
}

func getMessageBytesSend(pStateManager state.IStateManager, pR *http.Request) ([]byte, error) {
	switch pR.FormValue("method") {
	case http.MethodPost:
		rawMsg := strings.TrimSpace(pR.FormValue("input_message"))
		strMsg := utils.GetOnlyWritableCharacters(utils.ReplaceTextToEmoji(rawMsg))
		if strMsg == "" {
			return nil, errors.NewError("error: message is null")
		}
		return wrapText(strMsg), nil
	case http.MethodPut:
		filename, fileBytes, err := getUploadFile(pStateManager, pR)
		if err != nil {
			return nil, errors.WrapError(err, "error: upload file")
		}
		return wrapFile(filename, fileBytes), nil
	default:
		panic("got not supported method")
	}
}

func getUploadFile(pStateManager state.IStateManager, pR *http.Request) (string, []byte, error) {
	// Get handler for filename, size and headers
	file, handler, err := pR.FormFile("input_file")
	if err != nil {
		return "", nil, errors.NewError("error: receive file")
	}
	defer file.Close()

	if handler.Size == 0 {
		return "", nil, errors.NewError("error: file size is nil")
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", nil, errors.NewError("error: read file bytes")
	}

	return handler.Filename, fileBytes, nil
}

func trySendMessage(client client.IClient, recvPubKey, myPubKey asymmetric.IPubKey, msgBytes []byte, msgLimit uint64) error {
	// if the sender = receiver then there is no need to send a message to the network
	if myPubKey.GetAddress().ToString() == recvPubKey.GetAddress().ToString() {
		return nil
	}

	if uint64(len(msgBytes)) > msgLimit {
		return errors.NewError("error: len message > limit")
	}

	return client.BroadcastRequest(
		recvPubKey,
		request.NewRequest(http.MethodPost, pkg_settings.CTitlePattern, hlm_settings.CPushPath).
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody(msgBytes),
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
