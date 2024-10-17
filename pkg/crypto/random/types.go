package random

type IRandom interface {
	GetString(uint64) string
	GetBytes(uint64) []byte
	GetUint64() uint64
	GetBool() bool
}
