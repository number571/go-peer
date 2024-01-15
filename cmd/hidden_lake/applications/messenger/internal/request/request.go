package request

import (
	"net/http"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

func NewPushRequest(pPubKey asymmetric.IPubKey, pRequestID string, pBody []byte) hls_request.IRequest {
	address := pPubKey.GetHasher().ToBytes()
	return hls_request.NewRequest(http.MethodPost, hlm_settings.CTitlePattern, hlm_settings.CPushPath).
		WithHead(map[string]string{
			"Content-Type":                "application/json",
			hlm_settings.CHeaderSenderId:  encoding.HexEncode(address),
			hlm_settings.CHeaderRequestId: pRequestID,
		}).
		WithBody(pBody)
}
