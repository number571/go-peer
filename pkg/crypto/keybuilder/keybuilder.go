package keybuilder

import (
	"crypto/sha512"

	"golang.org/x/crypto/pbkdf2"
)

var (
	_ IKeyBuilder = &sKeyBuilder{}
)

type sKeyBuilder struct {
	fIterN uint64
	fSalt  []byte
}

func NewKeyBuilder(pIterN uint64, pSalt []byte) IKeyBuilder {
	return &sKeyBuilder{
		fIterN: pIterN,
		fSalt:  pSalt,
	}
}

func (p *sKeyBuilder) Build(pPassword string, pKeyLen uint64) []byte {
	return pbkdf2.Key(
		[]byte(pPassword),
		p.fSalt,
		int(p.fIterN),
		int(pKeyLen),
		sha512.New,
	)
}
