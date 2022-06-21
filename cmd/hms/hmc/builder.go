package hmc

import (
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/routing"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

var (
	_ IBuilder = &sBuiler{}
)

type sBuiler struct {
	client client.IClient
}

func NewBuiler(client client.IClient) IBuilder {
	return &sBuiler{
		client: client,
	}
}

func (builder *sBuiler) Size() *hms_settings.SSizeRequest {
	return &hms_settings.SSizeRequest{
		Receiver: builder.client.PubKey().Address().Bytes(),
	}
}

func (builder *sBuiler) Load(n uint64) *hms_settings.SLoadRequest {
	return &hms_settings.SLoadRequest{
		Receiver: builder.client.PubKey().Address().Bytes(),
		Index:    n,
	}
}

func (builder *sBuiler) Push(receiver asymmetric.IPubKey, msg []byte) *hms_settings.SPushRequest {
	encMsg, _ := builder.client.Encrypt(
		routing.NewRoute(builder.client.PubKey()),
		message.NewMessage([]byte(hms_settings.CTitlePattern), msg),
	)

	return &hms_settings.SPushRequest{
		Receiver: builder.client.PubKey().Address().Bytes(),
		Package:  encMsg.ToPackage().Bytes(),
	}
}
