package settings

const (
	MaskRout Key = iota + 1
	TimeWait
	TimePsdo
	SizePsdo
	SizeRtry
	SizeWork
	SizeConn
	SizePack
	SizeMapp
	SizeAkey
	SizeSkey
)

type (
	Key   = uint64
	Value = uint64
)
type Settings interface {
	Set(Key, Value)
	Get(Key) Value
}
