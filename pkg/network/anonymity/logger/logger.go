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

func NewLogger(service string) ILogger {
	if len(service) != 3 {
		return nil
	}
	return &sLogger{
		fService: service,
	}
}

func (l *sLogger) FmtLog(lType ILogType, msgHash []byte, proof uint64, pubKey asymmetric.IPubKey, netConn conn.IConn) string {
	conn := "127.0.0.1:"
	if netConn != nil {
		conn = netConn.Socket().RemoteAddr().String()
	}
	addr := make([]byte, hashing.CSHA256Size)
	if pubKey != nil {
		addr = pubKey.Address().Bytes()
	}
	hash := make([]byte, hashing.CSHA256Size)
	if msgHash != nil {
		hash = msgHash
	}
	return fmt.Sprintf(cLogTemplate, l.fService, lType, hash[:4], hash[28:], addr[:4], addr[28:], proof, conn)
}
