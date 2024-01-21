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
	GetOutInfo() io.Writer
	GetOutWarn() io.Writer
	GetOutErro() io.Writer
}
