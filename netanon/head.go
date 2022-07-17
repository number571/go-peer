package netanon

var (
	_ iHead = sHead(0)
)

// (32bit=A||32bit=B)
// A used for handle
// B used for action
type sHead uint64

func loadHead(n uint64) iHead {
	return sHead(n)
}

func joinHead(routes, actions uint32) iHead {
	return sHead((uint64(routes) << 32) | uint64(actions))
}

func (k sHead) Routes() uint32 {
	return uint32(k >> 32)
}

func (k sHead) Actions() uint32 {
	return uint32(k & 0x00000000FFFFFFFF)
}

func (k sHead) Uint64() uint64 {
	return uint64(k)
}
