package wrapper

type IWrapper interface {
	Get() interface{}
	Set(interface{}) IWrapper
}
