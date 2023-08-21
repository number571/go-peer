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
	CLogBaseMessageType     ILogType = "MTYPE"
	CLogBaseEnqueueRequest  ILogType = "ENQRQ"
	CLogBaseEnqueueResponse ILogType = "ENQRS"
	CLogBaseGetResponse     ILogType = "GETRS"

	// INFO
	CLogInfoExist           ILogType = "EXIST"
	CLogInfoUndecryptable   ILogType = "UNDEC"
	CLogInfoWithoutResponse ILogType = "WTHRS"

	// WARN
	CLogWarnMessageNull       ILogType = "MNULL"
	CLogWarnNotFriend         ILogType = "NTFRN"
	CLogWarnNotConnection     ILogType = "NTCON"
	CLogWarnUnknownRoute      ILogType = "UNKRT"
	CLogWarnIncorrectResponse ILogType = "ICRSP"

	// ERRO
	CLogErroDatabaseGet    ILogType = "DBGET"
	CLogErroDatabaseSet    ILogType = "DBSET"
	CLogErroEncryptPayload ILogType = "ENCPL"
)
