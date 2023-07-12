package logbuilder

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/network/conn"
)

const (
	cLogTemplate = "service=%s type=%s hash=%08X...%08X addr=%08X...%08X proof=%016d conn=%s"
)

type sLogger struct {
	fService string
	fHash    []byte
	fProof   uint64
	fPubKey  asymmetric.IPubKey
	fConn    conn.IConn
}

func NewLogBuilder(pService string) ILogBuilder {
	return &sLogger{
		fService: pService,
	}
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

func (p *sLogger) WithConn(pConn conn.IConn) ILogBuilder {
	p.fConn = pConn
	return p
}

func (p *sLogger) Get(pType ILogType) string {
	conn := "127.0.0.1:"
	if p.fConn != nil {
		conn = p.fConn.GetSocket().RemoteAddr().String()
	}
	addr := make([]byte, hashing.CSHA256Size)
	if p.fPubKey != nil {
		addr = p.fPubKey.GetAddress().ToBytes()
	}
	hash := make([]byte, hashing.CSHA256Size)
	if p.fHash != nil {
		hash = p.fHash
	}
	return fmt.Sprintf(cLogTemplate, p.fService, pType, hash[:4], hash[28:], addr[:4], addr[28:], p.fProof, conn)
}
