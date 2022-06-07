package hmc

import (
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

var (
	_ IBuilder = &sBuiler{}
)

type sBuiler struct {
	client local.IClient
}

func NewBuiler(client local.IClient) IBuilder {
	return &sBuiler{
		client: client,
	}
}

func (builder *sBuiler) Size() *hms_settings.SSizeRequest {
	pubBytes := builder.client.PubKey().Bytes()
	hashRecv := crypto.NewHasher(pubBytes).Bytes()

	return &hms_settings.SSizeRequest{
		Receiver: hashRecv,
	}
}

func (builder *sBuiler) Load(n uint64) *hms_settings.SLoadRequest {
	pubBytes := builder.client.PubKey().Bytes()
	hashRecv := crypto.NewHasher(pubBytes).Bytes()

	return &hms_settings.SLoadRequest{
		Receiver: hashRecv,
		Index:    n,
	}
}

func (builder *sBuiler) Push(receiver crypto.IPubKey, msg []byte) *hms_settings.SPushRequest {
	pubBytes := receiver.Bytes()
	hashRecv := crypto.NewHasher(pubBytes).Bytes()

	encMsg, _ := builder.client.Encrypt(
		local.NewRoute(receiver),
		local.NewMessage([]byte(hms_settings.CTitlePattern), msg),
	)

	return &hms_settings.SPushRequest{
		Receiver: hashRecv,
		Package:  encMsg.ToPackage().Bytes(),
	}
}
