package selector

import (
	"github.com/number571/go-peer/crypto/asymmetric"
)

type ISelector interface {
	Length() uint64
	Shuffle() ISelector
	Return(uint64) []asymmetric.IPubKey
}
