package anon

import anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"

var gLogMap = map[anon_logger.ILogType]string{
	0: "", // invalid log

	// default
	anon_logger.CLogBaseBroadcast:          "BRDCS",
	anon_logger.CLogBaseMessageType:        "MTYPE",
	anon_logger.CLogBaseEnqueueRequest:     "ENQRQ",
	anon_logger.CLogBaseEnqueueResponse:    "ENQRS",
	anon_logger.CLogBaseGetResponse:        "GETRS",
	anon_logger.CLogInfoExist:              "EXIST",
	anon_logger.CLogInfoUndecryptable:      "UNDEC",
	anon_logger.CLogInfoWithoutResponse:    "WTHRS",
	anon_logger.CLogWarnMessageNull:        "MNULL",
	anon_logger.CLogWarnNotFriend:          "NTFRN",
	anon_logger.CLogWarnNotConnection:      "NTCON",
	anon_logger.CLogWarnUnknownRoute:       "UNKRT",
	anon_logger.CLogWarnIncorrectResponse:  "ICRSP",
	anon_logger.CLogErroDatabaseGet:        "DBGET",
	anon_logger.CLogErroDatabaseSet:        "DBSET",
	anon_logger.CLogErroEncryptPayload:     "ENCPL",
	anon_logger.CLogErroInvalidNetworkMask: "IVNMS",

	// extend
	CLogBaseResponseModeFromService: "RSPMD",
	CLogInfoResponseFromService:     "RSPSR",
	CLogWarnRequestToService:        "RQTSR",
	CLogWarnUndefinedService:        "UNDSR",
	CLogErroLoadRequestType:         "LDRQT",
	CLogErroProxyRequestType:        "PXRQT",
}

const (
	// BASE
	CLogBaseResponseModeFromService anon_logger.ILogType = iota + anon_logger.CLogFinal + 1

	// INFO
	CLogInfoResponseFromService

	// WARN
	CLogWarnRequestToService
	CLogWarnUndefinedService

	// ERRO
	CLogErroLoadRequestType
	CLogErroProxyRequestType
)
