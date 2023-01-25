package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/settings"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/payload"
)

func HandleBroadcastAPI(db database.IKeyValueDB, connKeeper conn_keeper.IConnKeeper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		hashes, err := db.Hashes()
		if err != nil {
			response(w, pkg_settings.CErrorLoad, "failed: load hashes from db")
			return
		}

		for _, hash := range hashes {
			msg, err := db.Load(hash)
			if err != nil {
				continue // may be deleted
			}
			err = connKeeper.Network().Broadcast(payload.NewPayload(
				hlt_settings.CNetworkMask,
				msg.Bytes(),
			))
			if err != nil {
				// TODO: log
				continue
			}
		}

		response(w, pkg_settings.CErrorNone, pkg_settings.CTitlePattern)
	}
}
