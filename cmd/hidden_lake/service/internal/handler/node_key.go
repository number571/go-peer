package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNodeKeyAPI(pWrapper config.IWrapper, pNode anonymity.INode) http.HandlerFunc {
	keySize := pWrapper.GetConfig().GetKeySizeBits()
	ephPrivKey := asymmetric.NewRSAPrivKey(keySize)

	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			var vPrivKey pkg_settings.SPrivKey

			if err := json.NewDecoder(pR.Body).Decode(&vPrivKey); err != nil {
				api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}

			sessionKey := ephPrivKey.DecryptBytes(encoding.HexDecode(vPrivKey.FEncSessionKey))
			encPrivKey := encoding.HexDecode(vPrivKey.FEncPrivKey)
			privKeyBytes := symmetric.NewAESCipher(sessionKey).DecryptBytes(encPrivKey)

			privKey := asymmetric.LoadRSAPrivKey(privKeyBytes)
			if privKey == nil {
				api.Response(pW, http.StatusBadRequest, "failed: decode private key")
				return
			}

			if privKey.GetSize() != pWrapper.GetConfig().GetKeySizeBits() {
				api.Response(pW, http.StatusNotAcceptable, "failed: incorrect private key size")
				return
			}

			client := pkg_settings.InitClient(pWrapper.GetConfig(), privKey)
			pNode.GetMessageQueue().UpdateClient(client)
		}

		pubKey := pNode.GetMessageQueue().GetClient().GetPubKey().ToString()
		pubExp := ephPrivKey.GetPubKey().ToString()

		api.Response(pW, http.StatusOK, fmt.Sprintf("%s,%s", pubKey, pubExp))
	}
}
