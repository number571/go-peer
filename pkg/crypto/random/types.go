package random

type IPRNG interface {
	GetString(uint64) string
	GetBytes(uint64) []byte
	GetUint64() uint64
	GetBool() bool
}
