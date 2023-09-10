package handler

import (
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"

	pkg_client "github.com/number571/go-peer/pkg/client"
)

var (
	gMutex          sync.Mutex
	gMsgSizeBytes   = uint64(0)
	gKeySizeBytes   = uint64(0)
	gMsgLimitBase64 = uint64(0)
)

func getMessageLimit(pHlsClient client.IClient) (uint64, error) {
	gMutex.Lock()
	defer gMutex.Unlock()

	sett, err := pHlsClient.GetSettings()
	if err != nil {
		return 0, errors.WrapError(err, "get settings from HLS (message size)")
	}

	msgSize := sett.GetMessageSizeBytes()
	keySize := sett.GetKeySizeBits()

	if msgSize == gMsgSizeBytes && keySize == gKeySizeBytes {
		return gMsgLimitBase64, nil
	}

	randClient := pkg_client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: msgSize,
			FWorkSizeBits:     1, // does not affect the size
		}),
		asymmetric.NewRSAPrivKey(keySize),
	)

	// overhead base64 format: https://ru.wikipedia.org/wiki/Base64
	msgLimitBytes := randClient.GetMessageLimit()
	msgLimitBase64 := msgLimitBytes - (msgLimitBytes / 4)

	gMsgSizeBytes = msgSize
	gKeySizeBytes = keySize
	gMsgLimitBase64 = msgLimitBase64

	return msgLimitBase64, nil
}
