package encoding

import (
	"encoding/json"

	"github.com/number571/go-peer/pkg/utils"
)

func SerializeJSON(pData interface{}) []byte {
	res, _ := json.Marshal(pData)
	return res
}

func DeserializeJSON(pData []byte, pRes interface{}) error {
	if err := json.Unmarshal(pData, pRes); err != nil {
		return utils.MergeErrors(ErrDeserialize, err)
	}
	return nil
}
