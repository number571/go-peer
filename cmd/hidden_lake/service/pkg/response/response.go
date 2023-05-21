package response

import "github.com/number571/go-peer/pkg/encoding"

var (
	_ IResponse = &sResponse{}
)

type sResponse struct {
	FCode int    `json:"code"`
	FBody []byte `json:"body"`
}

func NewResponse(pCode int, pBody []byte) IResponse {
	return &sResponse{
		FCode: pCode,
		FBody: pBody,
	}
}

func LoadResponse(pBytes []byte) (IResponse, error) {
	response := new(sResponse)
	err := encoding.Deserialize(pBytes, response)
	return response, err
}

func (p *sResponse) ToBytes() []byte {
	return encoding.Serialize(p)
}

func (p *sResponse) GetCode() int {
	return p.FCode
}

func (p *sResponse) GetBody() []byte {
	return p.FBody
}
