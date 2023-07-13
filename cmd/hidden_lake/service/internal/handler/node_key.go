package handler

import (
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNodeKeyAPI(pWrapper config.IWrapper, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			// next
		case http.MethodPost:
			privKeyBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				api.Response(pW, http.StatusConflict, "failed: read private key bytes")
				return
			}

			privKey := asymmetric.LoadRSAPrivKey(string(privKeyBytes))
			if privKey == nil {
				api.Response(pW, http.StatusBadRequest, "failed: decode private key")
				return
			}

			if privKey.GetSize() != pkg_settings.CAKeySize {
				api.Response(pW, http.StatusNotAcceptable, "failed: incorrect private key size")
				return
			}

			client := pkg_settings.InitClient(pWrapper.GetConfig(), privKey)
			pNode.GetMessageQueue().UpdateClient(client)
		}

		api.Response(pW, http.StatusOK, pNode.GetMessageQueue().GetClient().GetPubKey().ToString())
	}
}
