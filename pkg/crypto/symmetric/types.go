package symmetric

import (
	"github.com/number571/go-peer/pkg/crypto"
	"github.com/number571/go-peer/pkg/types"
)

type ICipher interface {
	crypto.IEncrypter
	crypto.IDecrypter
	types.IConverter
	types.IParameter
}
