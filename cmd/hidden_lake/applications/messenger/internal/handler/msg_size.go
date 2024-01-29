package handler

import (
	"errors"
	"fmt"

	hlm_client "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/internal/base64"
)

var (
	gReqSize = uint64(len(hlm_client.NewBuilder().PushMessage(
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
		[]byte{},
	).GetBody()))
)

func getMessageLimit(pHlsClient client.IClient) (uint64, error) {
	sett, err := pHlsClient.GetSettings()
	if err != nil {
		return 0, fmt.Errorf("get settings from HLS (message size): %w", err)
	}

	msgLimitOrig := sett.GetLimitMessageSizeBytes()
	if gReqSize >= msgLimitOrig {
		return 0, errors.New("push message size >= limit message size")
	}

	msgLimitBytes := msgLimitOrig - gReqSize
	return base64.GetSizeInBase64(msgLimitBytes)
}
