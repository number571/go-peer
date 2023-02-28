package logger

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/conn"
)

type ILogger interface {
	GetFmtLog(ILogType, []byte, uint64, asymmetric.IPubKey, conn.IConn) string
}

type ILogType string

const (
	// Base
	CLogBaseBroadcast   ILogType = "BRDCS"
	CLogBaseEnqueueResp ILogType = "ENRSP"

	// INFO
	CLogInfoExist         ILogType = "EXIST"
	CLogInfoUnencryptable ILogType = "UNENC"
	CLogInfoAction        ILogType = "ACTON"
	CLogInfoWithoutResp   ILogType = "WHRSP"

	// WARN
	CLogWarnMessageNull  ILogType = "MNULL"
	CLogWarnNotFriend    ILogType = "NTFRN"
	CLogWarnUnknownRoute ILogType = "UNKRT"

	// ERRO
	CLogErroMiddleware  ILogType = "MDLWR"
	CLogErroDatabaseSet ILogType = "DBSET"
)
