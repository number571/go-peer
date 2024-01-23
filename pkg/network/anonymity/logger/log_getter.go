package logger

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/conn"
)

var (
	_ ILogGetter = &sLogGetter{}
)

type sLogGetter struct {
	fLogBuilder *sLogBuilder
}

func wrapLogBuilder(pLogBuilder ILogBuilder) ILogGetter {
	return &sLogGetter{
		fLogBuilder: pLogBuilder.(*sLogBuilder),
	}
}

func (p *sLogGetter) GetService() string {
	return p.fLogBuilder.fService
}

func (p *sLogGetter) GetType() ILogType {
	return p.fLogBuilder.fType
}

func (p *sLogGetter) GetConn() conn.IConn {
	return p.fLogBuilder.fConn
}

func (p *sLogGetter) GetHash() []byte {
	return p.fLogBuilder.fHash
}

func (p *sLogGetter) GetSize() uint64 {
	return p.fLogBuilder.fSize
}

func (p *sLogGetter) GetPubKey() asymmetric.IPubKey {
	return p.fLogBuilder.fPubKey
}

func (p *sLogGetter) GetProof() uint64 {
	return p.fLogBuilder.fProof
}
