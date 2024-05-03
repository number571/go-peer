package response

import (
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ IResponse = &SResponse{}
)

type SResponse struct {
	SResponseBlock
	FBody []byte `json:"body"`
}

type SResponseBlock struct {
	FCode int               `json:"code"`
	FHead map[string]string `json:"head"`
}

func NewResponse(pCode int) IResponse {
	return &SResponse{
		SResponseBlock: SResponseBlock{
			FCode: pCode,
		},
	}
}

func LoadResponse(pData interface{}) (IResponse, error) {
	var response = new(SResponse)
	switch x := pData.(type) {
	case []byte:
		bytesSlice, err := joiner.LoadBytesJoiner32(x)
		if err != nil || len(bytesSlice) != 2 {
			return nil, ErrLoadBytesJoiner
		}
		if err := encoding.DeserializeJSON(bytesSlice[0], response); err != nil {
			return nil, utils.MergeErrors(ErrDecodeResponse, err)
		}
		response.FBody = bytesSlice[1]
		return response, nil
	case string:
		if err := encoding.DeserializeJSON([]byte(x), response); err != nil {
			return nil, utils.MergeErrors(ErrDecodeResponse, err)
		}
		return response, nil
	default:
		return nil, ErrUnknownType
	}
}

func (p *SResponse) ToBytes() []byte {
	return joiner.NewBytesJoiner32([][]byte{
		encoding.SerializeJSON(p.SResponseBlock),
		p.FBody,
	})
}

func (p *SResponse) ToString() string {
	return string(encoding.SerializeJSON(p))
}

func (p *SResponse) WithHead(pHead map[string]string) IResponse {
	p.FHead = make(map[string]string)
	for k, v := range pHead {
		p.FHead[k] = v
	}
	return p
}

func (p *SResponse) WithBody(pBody []byte) IResponse {
	p.FBody = pBody
	return p
}

func (p *SResponse) GetCode() int {
	return p.FCode
}

func (p *SResponse) GetHead() map[string]string {
	headers := make(map[string]string)
	for k, v := range p.FHead {
		headers[k] = v
	}
	return headers
}

func (p *SResponse) GetBody() []byte {
	return p.FBody
}
