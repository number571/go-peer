package settings

import (
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func InitClient(pPrivKey asymmetric.IPrivKey) client.IClient {
	return client.NewClient(
		message.NewSettings(&message.SSettings{
			FWorkSize:    CWorkSize,
			FMessageSize: CMessageSize,
		}),
		pPrivKey,
	)
}
