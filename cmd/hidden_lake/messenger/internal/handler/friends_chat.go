package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/logger"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
)

type sChatMessage struct {
	FIsIncoming  bool
	FMessageInfo utils.SMessageInfo
}
type sChatAddress struct {
	FAliasName  string
	FPubKeyHash string
}
type sChatMessages struct {
	*state.STemplateState
	FAddress  sChatAddress
	FMessages []sChatMessage
}

func FriendsChatPage(pStateManager state.IStateManager, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/friends/chat" {
			NotFoundPage(pStateManager, pLogger)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogRedirect))
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		pR.ParseForm()
		pR.ParseMultipartForm(10 << 20) // default value

		aliasName := pR.URL.Query().Get("alias_name")
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("get_alias_name"))
			fmt.Fprint(pW, "alias name is null")
			return
		}

		db := pStateManager.GetWrapperDB().Get()
		if db == nil {
			pLogger.PushErro(logBuilder.WithMessage("get_database"))
			fmt.Fprint(pW, "error: database closed")
			return
		}

		client := pStateManager.GetClient()
		myPubKey, _, err := client.GetPubKey()
		if err != nil || !pStateManager.IsMyPubKey(myPubKey) {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			fmt.Fprint(pW, errors.WrapError(err, "error: read public key"))
			return
		}

		recvPubKey, err := getReceiverPubKey(client, myPubKey, aliasName)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_receiver"))
			fmt.Fprint(pW, errors.WrapError(err, "error: get receiver by public key"))
			return
		}

		rel := database.NewRelation(myPubKey, recvPubKey)

		switch pR.FormValue("method") {
		case http.MethodPost, http.MethodPut:
			msgBytes, err := getMessageBytes(pR)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("get_message"))
				fmt.Fprint(pW, errors.WrapError(err, "error: get message bytes"))
				return
			}

			if err := trySendMessage(client, aliasName, msgBytes); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("send_message"))
				fmt.Fprint(pW, errors.WrapError(err, "error: push message to network"))
				return
			}

			dbMsg := database.NewMessage(false, doMessageProcessor(msgBytes))
			if err := db.Push(rel, dbMsg); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("push_message"))
				fmt.Fprint(pW, errors.WrapError(err, "error: add message to database"))
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogRedirect))
			http.Redirect(pW, pR,
				fmt.Sprintf("/friends/chat?alias_name=%s", aliasName),
				http.StatusSeeOther)
			return
		}

		start := uint64(0)
		size := db.Size(rel)

		messagesCap := pStateManager.GetConfig().GetSettings().GetMessagesCapacity()
		if size > messagesCap {
			start = size - messagesCap
		}

		msgs, err := db.Load(rel, start, size)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("read_database"))
			fmt.Fprint(pW, errors.WrapError(err, "error: read database"))
			return
		}

		res := &sChatMessages{
			STemplateState: pStateManager.GetTemplate(),
			FAddress: sChatAddress{
				FAliasName:  aliasName,
				FPubKeyHash: recvPubKey.GetAddress().ToString(),
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

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		t.Execute(pW, res)
	}
}

func getMessageBytes(pR *http.Request) ([]byte, error) {
	switch pR.FormValue("method") {
	case http.MethodPost:
		strMsg := strings.TrimSpace(pR.FormValue("input_message"))
		if strMsg == "" {
			return nil, errors.NewError("error: message is null")
		}
		if utils.HasNotWritableCharacters(strMsg) {
			return nil, errors.NewError("error: message has not writable characters")
		}
		return wrapText(strMsg), nil
	case http.MethodPut:
		filename, fileBytes, err := getUploadFile(pR)
		if err != nil {
			return nil, errors.WrapError(err, "error: upload file")
		}
		return wrapFile(filename, fileBytes), nil
	default:
		panic("got not supported method")
	}
}

func getUploadFile(pR *http.Request) (string, []byte, error) {
	// Get handler for filename, size and headers
	file, handler, err := pR.FormFile("input_file")
	if err != nil {
		return "", nil, errors.WrapError(err, "error: receive file")
	}
	defer file.Close()

	if handler.Size == 0 {
		return "", nil, errors.NewError("error: file size is nil")
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", nil, errors.WrapError(err, "error: read file bytes")
	}

	return handler.Filename, fileBytes, nil
}

func trySendMessage(pClient client.IClient, pAliasName string, pMsgBytes []byte) error {
	msgLimit, err := getMessageLimit(pClient)
	if err != nil {
		return errors.WrapError(err, "error: try send message")
	}

	if uint64(len(pMsgBytes)) > msgLimit {
		return errors.NewError("error: len message > limit")
	}

	// if the sender = receiver then there is no need to send a message to the network
	if pAliasName == hlm_settings.CIamAliasName {
		return nil
	}

	return pClient.BroadcastRequest(
		pAliasName,
		request.NewRequest(http.MethodPost, hlm_settings.CTitlePattern, hlm_settings.CPushPath).
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody(pMsgBytes),
	)
}

func getReceiverPubKey(client client.IClient, myPubKey asymmetric.IPubKey, aliasName string) (asymmetric.IPubKey, error) {
	if aliasName == hlm_settings.CIamAliasName {
		return myPubKey, nil
	}

	friends, err := client.GetFriends()
	if err != nil {
		return nil, errors.WrapError(err, "error: read friends")
	}

	friendPubKey, ok := friends[aliasName]
	if !ok {
		return nil, errors.NewError("undefined public key by alias name")
	}

	return friendPubKey, nil
}
