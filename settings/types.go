package settings

type (
	Key   = uint64
	Value = uint64
)

type Settings interface {
	Set(Key, Value) Settings
	Get(Key) Value
}

const (
	MaskRout Key = iota + 1
	TimeWait
	TimePsdo
	SizeRtry
	SizeWork
	SizeConn
	SizePack
	SizeMapp
	SizeAkey
	SizeSkey
)
