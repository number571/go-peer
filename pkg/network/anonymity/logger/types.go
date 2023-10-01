package logger

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/conn"
)

type (
	ILogType uint8
)

const (
	CLogFinal = CLogErroEncryptPayload
)

const (
	// Base (can be in >1 more state)
	CLogBaseBroadcast ILogType = iota + 1
	CLogBaseMessageType
	CLogBaseEnqueueRequest
	CLogBaseEnqueueResponse
	CLogBaseGetResponse

	// INFO
	CLogInfoExist
	CLogInfoUndecryptable
	CLogInfoWithoutResponse

	// WARN
	CLogWarnMessageNull
	CLogWarnNotFriend
	CLogWarnNotConnection
	CLogWarnUnknownRoute
	CLogWarnIncorrectResponse

	// ERRO
	CLogErroDatabaseGet
	CLogErroDatabaseSet
	CLogErroEncryptPayload
)

type ILogGetter interface {
	GetService() string
	GetType() ILogType
	GetSize() uint64
	GetProof() uint64
	GetHash() []byte
	GetConn() conn.IConn
	GetPubKey() asymmetric.IPubKey
}

type ILogBuilder interface {
	Get() ILogGetter

	WithType(ILogType) ILogBuilder
	WithSize(int) ILogBuilder
	WithProof(uint64) ILogBuilder
	WithHash([]byte) ILogBuilder
	WithConn(conn.IConn) ILogBuilder
	WithPubKey(asymmetric.IPubKey) ILogBuilder
}
