package logger

import (
	"github.com/number571/go-peer/pkg/crypto/quantum"
	"github.com/number571/go-peer/pkg/network/conn"
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
	CLogInfoPassF2FOption
	CLogInfoWithoutResponse

	// WARN
	CLogWarnMessageType
	CLogWarnMessageNull
	CLogWarnPayloadNull
	CLogWarnNotFriend
	CLogWarnUnknownRoute
	CLogWarnIncorrectResponse

	// ERRO
	CLogErroDatabaseGet
	CLogErroDatabaseSet
)

type ILogBuilder interface {
	ILogGetterFactory

	WithType(ILogType) ILogBuilder
	WithSize(int) ILogBuilder
	WithProof(uint64) ILogBuilder
	WithHash([]byte) ILogBuilder
	WithConn(conn.IConn) ILogBuilder
	WithPubKey(quantum.ISignerPubKey) ILogBuilder
}

type ILogGetterFactory interface {
	Get() ILogGetter
}

type ILogGetter interface {
	GetService() string
	GetType() ILogType
	GetSize() uint64
	GetProof() uint64
	GetHash() []byte
	GetConn() conn.IConn
	GetPubKey() quantum.ISignerPubKey
}
