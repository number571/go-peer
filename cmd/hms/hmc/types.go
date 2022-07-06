package hmc

import (
	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/payload"
)

type IClient interface {
	Size() (uint64, error)
	Load(uint64) (message.IMessage, error)
	Push(asymmetric.IPubKey, payload.IPayload) error
}

type IBuilder interface {
	Size() *hms_settings.SSizeRequest
	Load(uint64) *hms_settings.SLoadRequest
	Push(asymmetric.IPubKey, payload.IPayload) *hms_settings.SPushRequest
}

type IRequester interface {
	Size(*hms_settings.SSizeRequest) (uint64, error)
	Load(*hms_settings.SLoadRequest) (message.IMessage, error)
	Push(*hms_settings.SPushRequest) error
}
