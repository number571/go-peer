package anonymity

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/network/conn"
)

type iLogType string

const (
	// INFO
	cLogInfoExist         iLogType = "EXIST"
	cLogInfoUnencryptable iLogType = "UNENC"
	cLogInfoAction        iLogType = "ACTON"
	cLogInfoWithoutResp   iLogType = "WHRSP"
	cLogInfoEnqueueResp   iLogType = "ENRSP"
	cLogInfoBroadcast     iLogType = "BRDCS"

	// WARN
	cLogWarnMessageNull  iLogType = "MNULL"
	cLogWarnNotFriend    iLogType = "NTFRN"
	cLogWarnUnknownRoute iLogType = "UNKRT"

	// ERRO
	cLogErroDatabaseSet iLogType = "DBSET"
)

const (
	cLogTemplate = "type=%05s hash=%08X...%08X addr=%08X...%08X proof=%016d conn=%s"
)

func fmtLog(lType iLogType, msgHash []byte, proof uint64, pubKey asymmetric.IPubKey, netConn conn.IConn) string {
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
	return fmt.Sprintf(cLogTemplate, lType, hash[:4], hash[28:], addr[:4], addr[28:], proof, conn)
}
