package encoding

import (
	"github.com/number571/go-peer/pkg/utils"
	yaml "gopkg.in/yaml.v2"
)

func SerializeYAML(pData interface{}) []byte {
	res, _ := yaml.Marshal(pData)
	return res
}

func DeserializeYAML(pData []byte, pRes interface{}) error {
	if err := yaml.Unmarshal(pData, pRes); err != nil {
		return utils.MergeErrors(ErrDeserialize, err)
	}
	return nil
}
