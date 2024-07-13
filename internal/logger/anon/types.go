package anon

import anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"

var gLogMap = map[anon_logger.ILogType]string{
	0: "", // invalid log

	// default
	anon_logger.CLogBaseBroadcast:         "BRDCS",
	anon_logger.CLogBaseEnqueueRequest:    "ENQRQ",
	anon_logger.CLogBaseEnqueueResponse:   "ENQRS",
	anon_logger.CLogBaseGetResponse:       "GETRS",
	anon_logger.CLogInfoExist:             "EXIST",
	anon_logger.CLogInfoUndecryptable:     "UNDEC",
	anon_logger.CLogInfoPassF2FOption:     "PF2FO",
	anon_logger.CLogInfoWithoutResponse:   "WTHRS",
	anon_logger.CLogWarnMessageType:       "MTYPE",
	anon_logger.CLogWarnMessageNull:       "MNULL",
	anon_logger.CLogWarnPayloadNull:       "PNULL",
	anon_logger.CLogWarnNotFriend:         "NTFRN",
	anon_logger.CLogWarnUnknownRoute:      "UNKRT",
	anon_logger.CLogWarnIncorrectResponse: "ICRSP",
	anon_logger.CLogErroDatabaseGet:       "DBGET",
	anon_logger.CLogErroDatabaseSet:       "DBSET",

	// extend
	CLogBaseResponseModeFromService: "RSPMD",
	CLogBaseAppendNewFriend:         "APNFR",
	CLogInfoResponseFromService:     "RSPSR",
	CLogInfoRequestIDAlreadyExist:   "RQIDE",
	CLogWarnRequestToService:        "RQTSR",
	CLogWarnUndefinedService:        "UNDSR",
	CLogWarnInvalidRequestID:        "INRID",
	CLogErroLoadRequestType:         "LDRQT",
	CLogErroPushDatabaseType:        "PSHDB",
	CLogErroProxyRequestType:        "PXRQT",
}

const (
	// BASE
	CLogBaseResponseModeFromService anon_logger.ILogType = iota + anon_logger.CLogFinal + 1
	CLogBaseAppendNewFriend

	// INFO
	CLogInfoResponseFromService
	CLogInfoRequestIDAlreadyExist

	// WARN
	CLogWarnRequestToService
	CLogWarnUndefinedService
	CLogWarnInvalidRequestID

	// ERRO
	CLogErroLoadRequestType
	CLogErroPushDatabaseType
	CLogErroProxyRequestType
)
