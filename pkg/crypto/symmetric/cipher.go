package symmetric

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/number571/go-peer/pkg/encoding"
)

type iEncMode int

const (
	modeCFB iEncMode = iota + 1
	modeGCM
)

var (
	_ ICipher = &sAESCipher{}
)

const (
	CCipherBlockSize = aes.BlockSize
	CCipherKeySize   = 32
)

type sAESCipher struct {
	fMode  iEncMode
	fKey   []byte
	fBlock cipher.Block
}

func (p *sAESCipher) EncryptBytes(pMsg []byte) []byte {
	switch p.fMode {
	case modeCFB:
		return p.encryptBytesCFB(pMsg)
	case modeGCM:
		return p.encryptBytesGCM(pMsg)
	default:
		panic("encryption mode undefined")
	}
}

func (p *sAESCipher) DecryptBytes(pMsg []byte) []byte {
	switch p.fMode {
	case modeCFB:
		return p.decryptBytesCFB(pMsg)
	case modeGCM:
		return p.decryptBytesGCM(pMsg)
	default:
		panic("encryption mode undefined")
	}
}

func (p *sAESCipher) ToBytes() []byte {
	return p.fKey
}

func (p *sAESCipher) ToString() string {
	return encoding.HexEncode(p.fKey)
}
