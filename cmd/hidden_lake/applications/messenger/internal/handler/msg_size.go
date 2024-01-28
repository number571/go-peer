package handler

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"

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

	sett, err := pHlsClient.GetSettings() // TODO: append to settings GetMessageLimit result
	if err != nil {
		return 0, fmt.Errorf("get settings from HLS (message size): %w", err)
	}

	msgSize := sett.GetMessageSizeBytes()
	keySize := sett.GetKeySizeBits()

	if msgSize == gMsgSizeBytes && keySize == gKeySizeBytes {
		return gMsgLimitBase64, nil
	}

	randClient := pkg_client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: msgSize,
			FKeySizeBits:      keySize,
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
