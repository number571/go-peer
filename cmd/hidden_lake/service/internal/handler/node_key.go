package handler

import (
	"io"
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNodeKeyAPI(pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			api.Response(pW, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			privKeyBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				api.Response(pW, pkg_settings.CErrorRead, "failed: read private key bytes")
				return
			}

			privKey := asymmetric.LoadRSAPrivKey(string(privKeyBytes))
			if privKey == nil {
				api.Response(pW, pkg_settings.CErrorPrivKey, "failed: decode private key")
				return
			}

			pNode.GetMessageQueue().UpdateClient(pkg_settings.InitClient(privKey))
		}

		// Response for GET and POST
		api.Response(pW, pkg_settings.CErrorNone, pNode.GetMessageQueue().GetClient().GetPubKey().ToString())
	}
}
