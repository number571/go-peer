package logger

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/network/conn"
)

const (
	cLogTemplate = "service=%03s type=%05s hash=%08X...%08X addr=%08X...%08X proof=%016d conn=%s"
)

type sLogger struct {
	fService string
}

func NewLogger(pService string) ILogger {
	if len(pService) != 3 {
		return nil
	}
	return &sLogger{
		fService: pService,
	}
}

func (p *sLogger) GetFmtLog(pType ILogType, pMsgHash []byte, pProof uint64, pPubKey asymmetric.IPubKey, pNetConn conn.IConn) string {
	conn := "127.0.0.1:"
	if pNetConn != nil {
		conn = pNetConn.GetSocket().RemoteAddr().String()
	}
	addr := make([]byte, hashing.CSHA256Size)
	if pPubKey != nil {
		addr = pPubKey.GetAddress().ToBytes()
	}
	hash := make([]byte, hashing.CSHA256Size)
	if pMsgHash != nil {
		hash = pMsgHash
	}
	return fmt.Sprintf(cLogTemplate, p.fService, pType, hash[:4], hash[28:], addr[:4], addr[28:], pProof, conn)
}
