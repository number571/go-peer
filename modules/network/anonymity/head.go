package anonymity

var (
	_ iHead = sHead(0)
)

// (32bit=A||32bit=B)
// A used for action
// B used for route
type sHead uint64

func loadHead(n uint64) iHead {
	return sHead(n)
}

func joinHead(action, route uint32) iHead {
	return sHead((uint64(action) << 32) | uint64(route))
}

func (k sHead) GetRoute() uint32 {
	return uint32(k & 0x00000000FFFFFFFF)
}

func (k sHead) GetAction() uint32 {
	return uint32(k >> 32)
}

func (k sHead) Uint64() uint64 {
	return uint64(k)
}
