package logger

import "os"

type ISettings interface {
	GetStreamInfo() *os.File
	GetStreamWarn() *os.File
	GetStreamErro() *os.File
}

type ILogger interface {
	GetSettings() ISettings
	PushInfo(string)
	PushWarn(string)
	PushErro(string)
}

type ILogging interface {
	HasInfo() bool
	HasWarn() bool
	HasErro() bool
}
