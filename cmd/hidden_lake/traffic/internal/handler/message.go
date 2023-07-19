package handler

import (
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
)

func HandleMessageAPI(pWrapperDB database.IWrapperDB) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		database := pWrapperDB.Get()

		switch pR.Method {
		case http.MethodGet:
			query := pR.URL.Query()

			msg, err := database.Load(query.Get("hash"))
			if err != nil {
				api.Response(pW, http.StatusNotFound, "failed: load message")
				return
			}

			api.Response(pW, http.StatusOK, encoding.HexEncode(msg.ToBytes()))
			return
		case http.MethodPost:
			msgBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}

			msg := message.LoadMessage(
				message.NewSettings(&message.SSettings{
					FMessageSizeBytes: database.Settings().GetMessageSizeBytes(),
					FWorkSizeBits:     database.Settings().GetWorkSizeBits(),
				}),
				msgBytes,
			)
			if msg == nil {
				api.Response(pW, http.StatusBadRequest, "failed: decode message")
				return
			}

			if err := database.Push(msg); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: push message")
				return
			}

			api.Response(pW, http.StatusOK, "success: handle message")
			return
		}
	}
}
