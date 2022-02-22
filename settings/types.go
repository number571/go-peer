package settings

type (
	Key   = uint64
	Value = uint64
)
type Settings interface {
	Set(Key, Value)
	Get(Key) Value
}
