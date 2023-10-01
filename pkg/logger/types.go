package logger

import "os"

type (
	ILogArg  interface{}
	ILogFunc func(ILogArg) string
)

type ILogger interface {
	GetSettings() ISettings

	PushInfo(ILogArg)
	PushWarn(ILogArg)
	PushErro(ILogArg)
}

type ISettings interface {
	GetStreamInfo() *os.File
	GetStreamWarn() *os.File
	GetStreamErro() *os.File
}
