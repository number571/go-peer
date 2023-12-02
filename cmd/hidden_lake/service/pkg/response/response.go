package response

import (
	"fmt"

	"github.com/number571/go-peer/pkg/encoding"
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
	if err := encoding.DeserializeJSON(pBytes, response); err != nil {
		return nil, fmt.Errorf("load response: %w", err)
	}
	return response, nil
}

func (p *sResponse) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *sResponse) ToString() string {
	return string(p.ToBytes())
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
