package handler

import (
	"encoding/json"
	"net/http"

	hls_settings "github.com/number571/go-peer/cmd/hls/internal/settings"
	pkg_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNodeKeyAPI(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vPrivKey pkg_settings.SPrivKey

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		switch r.Method {
		case http.MethodPost:
			if err := json.NewDecoder(r.Body).Decode(&vPrivKey); err != nil {
				response(w, pkg_settings.CErrorDecode, "failed: decode request")
				return
			}

			privKey := asymmetric.LoadRSAPrivKey(vPrivKey.FPrivKey)
			if privKey == nil {
				response(w, pkg_settings.CErrorPrivKey, "failed: decode private key")
				return
			}

			node.Queue().UpdateClient(hls_settings.InitClient(privKey))
		}

		// Response for GET and POST
		response(w, pkg_settings.CErrorNone, node.Queue().Client().PubKey().String())
	}
}
