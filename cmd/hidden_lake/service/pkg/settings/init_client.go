package settings

import (
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func InitClient(pSett message.ISettings, pPrivKey asymmetric.IPrivKey) client.IClient {
	return client.NewClient(
		message.NewSettings(&message.SSettings{
			FWorkSize:    pSett.GetWorkSize(),
			FMessageSize: pSett.GetMessageSize(),
		}),
		pPrivKey,
	)
}
