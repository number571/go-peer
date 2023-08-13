package logbuilder

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/conn"
)

type ILogBuilder interface {
	Get(ILogType) string

	WithSize(int) ILogBuilder
	WithProof(uint64) ILogBuilder
	WithHash([]byte) ILogBuilder
	WithConn(conn.IConn) ILogBuilder
	WithPubKey(asymmetric.IPubKey) ILogBuilder
}

type ILogType string

const (
	// Base
	CLogBaseBroadcast       ILogType = "BRDCS"
	CLogBaseEnqueueResponse ILogType = "ENRSP"

	// INFO
	CLogInfoExist           ILogType = "EXIST"
	CLogInfoUndecryptable   ILogType = "UNDEC"
	CLogInfoAction          ILogType = "ACTON"
	CLogInfoWithoutResponse ILogType = "WHRSP"

	// WARN
	CLogWarnMessageNull       ILogType = "MNULL"
	CLogWarnNotFriend         ILogType = "NTFRN"
	CLogWarnUnknownRoute      ILogType = "UNKRT"
	CLogWarnOldResponse       ILogType = "LDRSP"
	CLogWarnIncorrectResponse ILogType = "ICRSP"

	// ERRO
	CLogErroMessageType ILogType = "MTYPE"
	CLogErroDatabaseGet ILogType = "DBGET"
	CLogErroDatabaseSet ILogType = "DBSET"
)
