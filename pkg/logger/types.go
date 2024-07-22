package logger

import "io"

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
	GetInfoWriter() io.Writer
	GetWarnWriter() io.Writer
	GetErroWriter() io.Writer
}
