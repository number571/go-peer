package std

type ILogging interface {
	HasInfo() bool
	HasWarn() bool
	HasErro() bool
}
