package logger

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
)

var (
	_ ILogBuilder = &sLogger{}
	_ ILogGetter  = &sLogger{}
)

type sLogger struct {
	fService string
	fType    ILogType
	fHash    []byte
	fProof   uint64
	fSize    uint64
	fPubKey  asymmetric.IPubKey
}

func NewLogBuilder(pService string) ILogBuilder {
	logger := &sLogger{
		fService: pService,
		fHash:    make([]byte, hashing.CHasherSize),
	}
	return logger
}

func (p *sLogger) GetService() string {
	return p.fService
}

func (p *sLogger) GetType() ILogType {
	return p.fType
}

func (p *sLogger) GetHash() []byte {
	return p.fHash
}

func (p *sLogger) GetSize() uint64 {
	return p.fSize
}

func (p *sLogger) GetPubKey() asymmetric.IPubKey {
	return p.fPubKey
}

func (p *sLogger) GetProof() uint64 {
	return p.fProof
}

func (p *sLogger) Build() ILogGetter {
	return p
}

func (p *sLogger) WithType(pType ILogType) ILogBuilder {
	p.fType = pType
	return p
}

func (p *sLogger) WithHash(pHash []byte) ILogBuilder {
	p.fHash = pHash
	return p
}

func (p *sLogger) WithProof(pProof uint64) ILogBuilder {
	p.fProof = pProof
	return p
}

func (p *sLogger) WithPubKey(pPubKey asymmetric.IPubKey) ILogBuilder {
	p.fPubKey = pPubKey
	return p
}

func (p *sLogger) WithSize(pSize int) ILogBuilder {
	p.fSize = uint64(pSize)
	return p
}
