package utils

import (
	"context"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_response "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	gRespSize = uint64(len(
		hls_response.NewResponse(200).
			WithHead(map[string]string{
				"Content-Type":                   api.CApplicationOctetStream,
				hls_settings.CHeaderResponseMode: hls_settings.CHeaderResponseModeON,
			}).
			WithBody([]byte{}).
			ToBytes(),
	))
)

func GetMessageLimit(pCtx context.Context, pHlsClient hls_client.IClient) (uint64, error) {
	sett, err := pHlsClient.GetSettings(pCtx)
	if err != nil {
		return 0, utils.MergeErrors(ErrGetSettingsHLS, err)
	}

	msgLimitOrig := sett.GetLimitMessageSizeBytes()
	if gRespSize >= msgLimitOrig {
		return 0, ErrMessageSizeGteLimit
	}

	return msgLimitOrig - gRespSize, nil
}
