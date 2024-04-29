package anonymity

type iHead interface {
	uint64() uint64
	getRoute() uint32
	getAction() iAction
}

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

func joinHead(pAction iAction, pRoute uint32) iHead {
	return sHead((uint64(pAction.(sAction)) << 32) | uint64(pRoute))
}

func (p sHead) getRoute() uint32 {
	return uint32(p & 0xFFFFFFFF)
}

func (p sHead) getAction() iAction {
	return sAction(p >> 32)
}

func (p sHead) uint64() uint64 {
	return uint64(p)
}
