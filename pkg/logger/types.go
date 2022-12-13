package logger

import "os"

type ILogger interface {
	Info(string)
	Warn(string)
	Erro(string)
}

type ISettings interface {
	GetStreamInfo() *os.File
	GetStreamWarn() *os.File
	GetStreamErro() *os.File
}
