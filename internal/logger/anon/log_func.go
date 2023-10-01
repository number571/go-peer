package anon

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/logger"

	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

const (
	cLogTemplate = "service=%s type=%s hash=%08X...%08X addr=%08X...%08X proof=%010d size=%dB conn=%s"
)

func GetLogFunc() logger.ILogFunc {
	return func(pLogArg logger.ILogArg) string {
		logBuilder, ok := pLogArg.(anon_logger.ILogBuilder)
		if !ok {
			panic("got invalid log arg")
		}

		logGetter := logBuilder.Get()
		logType := logGetter.GetType()
		if logType == 0 {
			panic("got invalid log type")
		}

		logStrType, ok := gLogMap[logType]
		if !ok {
			panic("value not found in log map")
		}

		return getLog(logStrType, logGetter)
	}
}

func getLog(logStrType string, pLogGetter anon_logger.ILogGetter) string {
	conn := "127.0.0.1:"
	if x := pLogGetter.GetConn(); x != nil {
		conn = x.GetSocket().RemoteAddr().String()
	}

	addr := make([]byte, hashing.CSHA256Size)
	if x := pLogGetter.GetPubKey(); x != nil {
		addr = x.GetAddress().ToBytes()
	}

	hash := make([]byte, hashing.CSHA256Size)
	if x := pLogGetter.GetHash(); x != nil {
		hash = x
	}

	return fmt.Sprintf(
		cLogTemplate,
		pLogGetter.GetService(),
		logStrType,
		hash[:4], hash[28:],
		addr[:4], addr[28:],
		pLogGetter.GetProof(),
		pLogGetter.GetSize(),
		conn,
	)
}
