package std

const (
	CLogInfo = "info"
	CLogWarn = "warn"
	CLogErro = "erro"
)

type ILogging interface {
	HasInfo() bool
	HasWarn() bool
	HasErro() bool
}
