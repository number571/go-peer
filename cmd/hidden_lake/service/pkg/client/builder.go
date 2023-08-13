package client

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) SetPrivKey(pEphPubKey asymmetric.IPubKey, pPrivKey asymmetric.IPrivKey) *pkg_settings.SPrivKey {
	sessionKey := random.NewStdPRNG().GetBytes(32)
	return &pkg_settings.SPrivKey{
		FSessionKey: encoding.HexEncode(pEphPubKey.EncryptBytes(sessionKey)),
		FPrivKey:    encoding.HexEncode(symmetric.NewAESCipher(sessionKey).EncryptBytes(pPrivKey.ToBytes())),
	}
}

func (p *sBuilder) Friend(pAliasName string, pPubKey asymmetric.IPubKey) *pkg_settings.SFriend {
	if pPubKey == nil {
		return &pkg_settings.SFriend{
			FAliasName: pAliasName,
		}
	}
	return &pkg_settings.SFriend{
		FAliasName: pAliasName,
		FPublicKey: pPubKey.ToString(),
	}
}

func (p *sBuilder) Request(pReceiver string, pReq request.IRequest) *pkg_settings.SRequest {
	return &pkg_settings.SRequest{
		FReceiver: pReceiver,
		FReqData:  pReq.ToString(),
	}
}

func (p *sBuilder) Message(pMsg message.IMessage) string {
	return pMsg.ToString()
}
