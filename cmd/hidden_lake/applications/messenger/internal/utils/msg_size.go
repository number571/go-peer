package utils

import (
	"context"

	hlm_client "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/internal/base64"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	// {"method":"POST","host":"hidden-lake-messenger","path":"/push","head":null,"body":""}
	gReqSize = uint64(len(
		hlm_client.NewBuilder().PushMessage([]byte{}).ToBytes(),
	))
)

func GetMessageLimit(pCtx context.Context, pHlsClient client.IClient) (uint64, error) {
	sett, err := pHlsClient.GetSettings(pCtx)
	if err != nil {
		return 0, utils.MergeErrors(ErrGetSettingsHLS, err)
	}

	msgLimitOrig := sett.GetLimitMessageSizeBytes()
	if gReqSize >= msgLimitOrig {
		return 0, ErrMessageSizeGtLimit
	}

	msgLimitBytes := msgLimitOrig - gReqSize
	size, err := base64.GetSizeInBase64(msgLimitBytes)
	if err != nil {
		return 0, utils.MergeErrors(ErrGetSizeInBase64, err)
	}

	return size, nil
}
