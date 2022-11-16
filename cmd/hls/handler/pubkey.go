package handler

import (
	"net/http"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/network/anonymity"
)

func HandlePubKeyAPI(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response(w, hls_settings.CErrorNone, node.Queue().Client().PubKey().String())
	}
}
