package handler

import (
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
	"github.com/number571/go-peer/modules/network/anonymity"
)

func HandlePubKeyAPI(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response(w, pkg_settings.CErrorNone, node.Queue().Client().PubKey().String())
	}
}
