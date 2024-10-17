package logger

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/conn"
)

var (
	_ ILogBuilder = &sLogBuilder{}
)

type sLogBuilder struct {
	fService string
	fType    ILogType
	fHash    []byte
	fProof   uint64
	fSize    uint64
	fPubKey  asymmetric.ISignPubKey
	fConn    conn.IConn

	fGetter ILogGetter
}

func NewLogBuilder(pService string) ILogBuilder {
	logBuilder := &sLogBuilder{
		fService: pService,
	}
	logBuilder.fGetter = wrapLogBuilder(logBuilder)
	return logBuilder
}

func (p *sLogBuilder) Get() ILogGetter {
	return p.fGetter
}

func (p *sLogBuilder) WithType(pType ILogType) ILogBuilder {
	p.fType = pType
	return p
}

func (p *sLogBuilder) WithHash(pHash []byte) ILogBuilder {
	p.fHash = pHash
	return p
}

func (p *sLogBuilder) WithProof(pProof uint64) ILogBuilder {
	p.fProof = pProof
	return p
}

func (p *sLogBuilder) WithPubKey(pPubKey asymmetric.ISignPubKey) ILogBuilder {
	p.fPubKey = pPubKey
	return p
}

func (p *sLogBuilder) WithConn(pConn conn.IConn) ILogBuilder {
	p.fConn = pConn
	return p
}

func (p *sLogBuilder) WithSize(pSize int) ILogBuilder {
	p.fSize = uint64(pSize)
	return p
}
