package request

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
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
	prng := random.NewStdPRNG()
	return &sRequest{
		FMethod: pMethod,
		FHost:   pHost,
		FPath:   pPath,
		FHead: map[string]string{
			settings.CHeaderRequestId: prng.GetString(settings.CRequestIDSize),
		},
	}
}

func LoadRequest(pData interface{}) (IRequest, error) {
	request := new(sRequest)
	switch x := pData.(type) {
	case []byte:
		if err := encoding.DeserializeJSON(x, request); err != nil {
			return nil, fmt.Errorf("load request: %w", err)
		}
		return request, nil
	case string:
		return LoadRequest([]byte(x))
	default:
		panic("type not supported")
	}
}

func (p *sRequest) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *sRequest) ToString() string {
	return string(p.ToBytes())
}

func (p *sRequest) WithHead(pHead map[string]string) IRequest {
	for k, v := range pHead {
		if k == settings.CHeaderRequestId {
			panic("an attempt to overwrite the header")
		}
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
