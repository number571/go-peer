package encoding

import (
	"encoding/json"
)

func SerializeJSON(pData interface{}) []byte {
	res, err := json.Marshal(pData)
	if err != nil {
		return nil
	}
	return res
}

func DeserializeJSON(pData []byte, pRes interface{}) error {
	return json.Unmarshal(pData, pRes)
}
