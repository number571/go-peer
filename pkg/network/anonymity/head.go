package anonymity

var (
	_ iHead = sHead(0)
)

// (32bit=A||32bit=B)
// A used for action
// B used for route
type sHead uint64

func loadHead(pN uint64) iHead {
	return sHead(pN)
}

func joinHead(pAction, pRoute uint32) iHead {
	return sHead((uint64(pAction) << 32) | uint64(pRoute))
}

func (p sHead) GetRoute() uint32 {
	return uint32(p & 0x00000000FFFFFFFF)
}

func (p sHead) GetAction() uint32 {
	return uint32(p >> 32)
}

func (p sHead) Uint64() uint64 {
	return uint64(p)
}
