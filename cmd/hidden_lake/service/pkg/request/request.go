package request

import (
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
)

var (
	_ IRequest = &sRequest{}
)

type sRequest struct {
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

func LoadRequest(pData interface{}) (IRequest, error) {
	request := new(sRequest)
	switch x := pData.(type) {
	case []byte:
		if err := encoding.Deserialize(x, request); err != nil {
			return nil, errors.WrapError(err, "load request")
		}
		return request, nil
	case string:
		return LoadRequest([]byte(x))
	default:
		panic("type not supported")
	}
}

func (p *sRequest) ToBytes() []byte {
	return encoding.Serialize(p, false)
}

func (p *sRequest) ToString() string {
	return string(encoding.Serialize(p, true))
}

func (p *sRequest) WithHead(pHead map[string]string) IRequest {
	p.FHead = make(map[string]string)
	for k, v := range pHead {
		p.FHead[k] = v
	}
	return p
}

func (p *sRequest) WithBody(pBody []byte) IRequest {
	p.FBody = pBody
	return p
}

func (p *sRequest) GetHost() string {
	return p.FHost
}

func (p *sRequest) GetPath() string {
	return p.FPath
}

func (p *sRequest) GetMethod() string {
	return p.FMethod
}

func (p *sRequest) GetHead() map[string]string {
	headers := make(map[string]string)
	for k, v := range p.FHead {
		headers[k] = v
	}
	return headers
}

func (p *sRequest) GetBody() []byte {
	return p.FBody
}
