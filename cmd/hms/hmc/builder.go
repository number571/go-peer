package hmc

import (
	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/payload"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
	fClient client.IClient
}

func NewBuilder(client client.IClient) IBuilder {
	return &sBuilder{
		fClient: client,
	}
}

func (builder *sBuilder) Size() *hms_settings.SSizeRequest {
	return &hms_settings.SSizeRequest{
		FReceiver: builder.fClient.PubKey().Address().Bytes(),
	}
}

func (builder *sBuilder) Load(n uint64) *hms_settings.SLoadRequest {
	return &hms_settings.SLoadRequest{
		FReceiver: builder.fClient.PubKey().Address().Bytes(),
		FIndex:    n,
	}
}

func (builder *sBuilder) Push(receiver asymmetric.IPubKey, pl payload.IPayload) *hms_settings.SPushRequest {
	encMsg, err := builder.fClient.Encrypt(builder.fClient.PubKey(), pl)
	if err != nil {
		panic(err)
	}

	return &hms_settings.SPushRequest{
		FReceiver: builder.fClient.PubKey().Address().Bytes(),
		FPackage:  encMsg.Bytes(),
	}
}
