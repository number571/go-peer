package std

import (
	"github.com/number571/go-peer/pkg/logger"
)

func GetLogFunc() logger.ILogFunc {
	return func(pLogArg logger.ILogArg) string {
		log, ok := pLogArg.(string)
		if !ok {
			panic("got invalid log arg")
		}
		return log
	}
}
