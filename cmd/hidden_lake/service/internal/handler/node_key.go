package handler

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"

	http_logger "github.com/number571/go-peer/internal/logger/http"
)

func HandleNodeKeyAPI(pWrapper config.IWrapper, pLogger logger.ILogger, pNode anonymity.INode, pEphPrivKey asymmetric.IPrivKey) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pLogger.PushWarn(httpLogger.Get(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			var vPrivKey pkg_settings.SPrivKey

			if err := json.NewDecoder(pR.Body).Decode(&vPrivKey); err != nil {
				pLogger.PushWarn(httpLogger.Get(http_logger.CLogDecodeBody))
				api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}

			privKey := getPrivKey(pEphPrivKey, vPrivKey)
			if privKey == nil {
				pLogger.PushWarn(httpLogger.Get("decode_key"))
				api.Response(pW, http.StatusBadRequest, "failed: decode private key")
				return
			}

			if privKey.GetSize() != pWrapper.GetConfig().GetKeySizeBits() {
				pLogger.PushWarn(httpLogger.Get("key_size"))
				api.Response(pW, http.StatusNotAcceptable, "failed: incorrect private key size")
				return
			}

			client := pkg_settings.InitClient(pWrapper.GetConfig(), privKey)
			pNode.GetMessageQueue().UpdateClient(client)

			pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: update private key")
			return
		}

		pubKeys := []string{
			pNode.GetMessageQueue().GetClient().GetPubKey().ToString(),
			pEphPrivKey.GetPubKey().ToString(),
		}

		pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, pubKeys)
	}
}

func getPrivKey(pEphPrivKey asymmetric.IPrivKey, pPrivKey pkg_settings.SPrivKey) asymmetric.IPrivKey {
	if pPrivKey.FSessionKey == "" {
		return asymmetric.LoadRSAPrivKey(pPrivKey.FPrivKey) // string
	}
	sessionKey := pEphPrivKey.DecryptBytes(encoding.HexDecode(pPrivKey.FSessionKey))
	encPrivKey := encoding.HexDecode(pPrivKey.FPrivKey)
	decPrivKey := symmetric.NewAESCipher(sessionKey).DecryptBytes(encPrivKey)
	return asymmetric.LoadRSAPrivKey(decPrivKey) // bytes
}
