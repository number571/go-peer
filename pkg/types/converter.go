package types

import (
	"encoding/json"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IConverter = &sConverter{}
)

type sConverter struct {
	fAsStr bool
	fBytes []byte
}

func NewConverter(pData interface{}) IConverter {
	switch x := pData.(type) {
	case []byte:
		return &sConverter{fBytes: x}
	case string:
		return &sConverter{fAsStr: true, fBytes: []byte(x)}
	default:
		v, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}
		return &sConverter{fAsStr: true, fBytes: v}
	}
}

func (p *sConverter) ToBytes() []byte {
	return p.fBytes
}

func (p *sConverter) ToString() string {
	if p.fAsStr {
		return string(p.fBytes)
	}
	return encoding.HexEncode(p.fBytes)
}
