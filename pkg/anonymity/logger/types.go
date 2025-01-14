package logger

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type (
	ILogType uint8
)

const (
	CLogFinal = CLogErroDatabaseSet
)

const (
	// Base (can be in >1 more state)
	CLogBaseBroadcast ILogType = iota + 1
	CLogBaseEnqueueRequest
	CLogBaseEnqueueResponse
	CLogBaseGetResponse

	// INFO
	CLogInfoExist
	CLogInfoUndecryptable
	CLogInfoWithoutResponse

	// WARN
	CLogWarnMessageNull
	CLogWarnPayloadNull
	CLogWarnUnknownRoute
	CLogWarnIncorrectResponse

	// ERRO
	CLogErroDatabaseGet
	CLogErroDatabaseSet
)

type ILogBuilder interface {
	Build() ILogGetter

	WithType(ILogType) ILogBuilder
	WithSize(int) ILogBuilder
	WithProof(uint64) ILogBuilder
	WithHash([]byte) ILogBuilder
	WithConn(string) ILogBuilder
	WithPubKey(asymmetric.IPubKey) ILogBuilder
}

type ILogGetter interface {
	GetService() string
	GetType() ILogType
	GetSize() uint64
	GetProof() uint64
	GetHash() []byte
	GetConn() string
	GetPubKey() asymmetric.IPubKey
}
