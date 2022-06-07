package hmc

import (
	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
)

type IClient interface {
	Size() (uint64, error)
	Load(uint64) (local.IMessage, error)
	Push(crypto.IPubKey, []byte) error
}

type IBuilder interface {
	Size() *hms_settings.SSizeRequest
	Load(uint64) *hms_settings.SLoadRequest
	Push(crypto.IPubKey, []byte) *hms_settings.SPushRequest
}

type IRequester interface {
	Size(*hms_settings.SSizeRequest) (uint64, error)
	Load(*hms_settings.SLoadRequest) (local.IMessage, error)
	Push(*hms_settings.SPushRequest) error
}
