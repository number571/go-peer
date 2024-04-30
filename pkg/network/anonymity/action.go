package anonymity

const (
	cAction32bitMask = sAction(1 << 31)
)

type iAction interface {
	uint31() uint32
	isRequest() bool
	setType(bool) iAction
}

var (
	_ iAction = sAction(0)
)

// (1bit=A||31bit=B)
// A = used as req=0/rsp=1
// B = used as action
type sAction uint32

func (p sAction) setType(isRequest bool) iAction {
	if isRequest {
		return p & ^cAction32bitMask
	}
	return p | cAction32bitMask
}

func (p sAction) isRequest() bool {
	f := p & cAction32bitMask
	return f == 0
}

func (p sAction) uint31() uint32 {
	return uint32(p) & uint32(^cAction32bitMask)
}
