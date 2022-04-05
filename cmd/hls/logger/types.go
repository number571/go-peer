package logger

type ILogger interface {
	Info(string)
	Warning(string)
	Error(string)
}
