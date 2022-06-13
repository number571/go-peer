package hmc

import (
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"
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
	pubBytes := builder.client.PubKey().Bytes()
	hashRecv := hashing.NewSHA256Hasher(pubBytes).Bytes()

	return &hms_settings.SSizeRequest{
		Receiver: hashRecv,
	}
}

func (builder *sBuiler) Load(n uint64) *hms_settings.SLoadRequest {
	pubBytes := builder.client.PubKey().Bytes()
	hashRecv := hashing.NewSHA256Hasher(pubBytes).Bytes()

	return &hms_settings.SLoadRequest{
		Receiver: hashRecv,
		Index:    n,
	}
}

func (builder *sBuiler) Push(receiver asymmetric.IPubKey, msg []byte) *hms_settings.SPushRequest {
	pubBytes := receiver.Bytes()
	hashRecv := hashing.NewSHA256Hasher(pubBytes).Bytes()

	encMsg, _ := builder.client.Encrypt(
		routing.NewRoute(receiver),
		message.NewMessage([]byte(hms_settings.CTitlePattern), msg),
	)

	return &hms_settings.SPushRequest{
		Receiver: hashRecv,
		Package:  encMsg.ToPackage().Bytes(),
	}
}
