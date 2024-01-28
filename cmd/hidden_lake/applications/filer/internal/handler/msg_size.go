package handler

import (
	"errors"
	"fmt"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_response "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/internal/base64"
)

var (
	// {"code":200,"head":{"Content-Type":"application/octet-stream","Hl-Service-Response-Mode":"on"},"body":""}
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

func getMessageLimit(pHlsClient hls_client.IClient) (uint64, error) {
	sett, err := pHlsClient.GetSettings()
	if err != nil {
		return 0, fmt.Errorf("get settings from HLS (message size): %w", err)
	}

	msgLimitOrig := sett.GetLimitMessageSizeBytes()
	if gRespSize >= msgLimitOrig {
		return 0, errors.New("response size >= limit message size")
	}

	msgLimitBytes := msgLimitOrig - gRespSize
	return base64.GetSizeInBase64(msgLimitBytes), nil
}
