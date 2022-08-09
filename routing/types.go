package routing

import (
	"github.com/number571/go-peer/crypto/asymmetric"
)

type IRoute interface {
	WithRedirects(asymmetric.IPrivKey, []asymmetric.IPubKey) IRoute

	Receiver() asymmetric.IPubKey
	PSender() asymmetric.IPrivKey
	List() []asymmetric.IPubKey
}
