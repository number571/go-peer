package logger

type ILogging interface {
	Info() bool
	Warn() bool
	Erro() bool
}
