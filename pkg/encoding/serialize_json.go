package encoding

import (
	"encoding/json"

	"github.com/number571/go-peer/pkg/utils"
)

func SerializeJSON(pData interface{}) []byte {
	res, err := json.Marshal(pData)
	if err != nil {
		return nil
	}
	return res
}

func DeserializeJSON(pData []byte, pRes interface{}) error {
	if err := json.Unmarshal(pData, pRes); err != nil {
		return utils.MergeErrors(ErrDeserialize, err)
	}
	return nil
}
