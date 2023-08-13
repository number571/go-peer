package logger

type ILogging interface {
	HasInfo() bool
	HasWarn() bool
	HasErro() bool
}
