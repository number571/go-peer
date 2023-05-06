package request

import (
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IRequest = &sRequest{}
)

type sRequest struct {
	fMutex  sync.Mutex
	FMethod string            `json:"method"`
	FHost   string            `json:"host"`
	FPath   string            `json:"path"`
	FHead   map[string]string `json:"head"`
	FBody   []byte            `json:"body"`
}

func NewRequest(pMethod, pHost, pPath string) IRequest {
	return &sRequest{
		FMethod: pMethod,
		FHost:   pHost,
		FPath:   pPath,
	}
}

func LoadRequest(pData []byte) (IRequest, error) {
	request := new(sRequest)
	err := encoding.Deserialize(pData, request)
	return request, err
}

func (p *sRequest) ToBytes() []byte {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return encoding.Serialize(p)
}

func (p *sRequest) WithHead(pHead map[string]string) IRequest {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.FHead = make(map[string]string)
	for k, v := range pHead {
		p.FHead[k] = v
	}
	return p
}

func (p *sRequest) WithBody(pBody []byte) IRequest {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.FBody = pBody
	return p
}

func (p *sRequest) Host() string {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.FHost
}

func (p *sRequest) Path() string {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.FPath
}

func (p *sRequest) Method() string {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.FMethod
}

func (p *sRequest) Head() map[string]string {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	headers := make(map[string]string)
	for k, v := range p.FHead {
		headers[k] = v
	}

	return headers
}

func (p *sRequest) Body() []byte {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.FBody
}
