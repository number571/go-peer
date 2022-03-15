package settings

type (
	Key   = uint64
	Value = uint64
)
type ISettings interface {
	Set(Key, Value) ISettings
	Get(Key) Value
}
