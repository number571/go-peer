package random

type IPRNG interface {
	String(uint64) string
	Bytes(uint64) []byte
	Uint64() uint64
	Bool() bool
}
