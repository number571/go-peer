package settings

type iSettings interface {
	Set(uint64, interface{}) iSettings
	Get(uint64) interface{}
}
