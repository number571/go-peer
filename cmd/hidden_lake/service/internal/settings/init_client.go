package settings

import (
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func InitClient(privKey asymmetric.IPrivKey) client.IClient {
	return client.NewClient(
		client.NewSettings(&client.SSettings{
			FWorkSize:    pkg_settings.CWorkSize,
			FMessageSize: pkg_settings.CMessageSize,
		}),
		privKey,
	)
}
