package keybuilder

import (
	"crypto/sha256"

	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"golang.org/x/crypto/pbkdf2"
)

var (
	_ IKeyBuilder = &sKeyBuilder{}
)

type sKeyBuilder struct {
	fSalt []byte
	fBits uint64
}

func NewKeyBuilder(pBits uint64, pSalt []byte) IKeyBuilder {
	return &sKeyBuilder{
		fBits: pBits,
		fSalt: pSalt,
	}
}

func (p *sKeyBuilder) Build(pPassword string) []byte {
	return pbkdf2.Key(
		[]byte(pPassword),
		p.fSalt,
		(1 << p.fBits),
		int(symmetric.CAESKeySize),
		sha256.New,
	)
}
