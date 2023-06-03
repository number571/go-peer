package response

import (
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
)

var (
	_ IResponse = &sResponse{}
)

type sResponse struct {
	FCode int               `json:"code"`
	FHead map[string]string `json:"head"`
	FBody []byte            `json:"body"`
}

func NewResponse(pCode int) IResponse {
	return &sResponse{
		FCode: pCode,
	}
}

func LoadResponse(pBytes []byte) (IResponse, error) {
	response := new(sResponse)
	if err := encoding.Deserialize(pBytes, response); err != nil {
		return nil, errors.WrapError(err, "load response")
	}
	return response, nil
}

func (p *sResponse) ToBytes() []byte {
	return encoding.Serialize(p, false)
}

func (p *sResponse) WithHead(pHead map[string]string) IResponse {
	p.FHead = make(map[string]string)
	for k, v := range pHead {
		p.FHead[k] = v
	}
	return p
}

func (p *sResponse) WithBody(pBody []byte) IResponse {
	p.FBody = pBody
	return p
}

func (p *sResponse) GetCode() int {
	return p.FCode
}

func (p *sResponse) GetHead() map[string]string {
	headers := make(map[string]string)
	for k, v := range p.FHead {
		headers[k] = v
	}
	return headers
}

func (p *sResponse) GetBody() []byte {
	return p.FBody
}
