package encoding

import (
	"encoding/json"
	"errors"
)

func SerializeJSON(pData interface{}) []byte {
	res, _ := json.Marshal(pData)
	return res
}

func DeserializeJSON(pData []byte, pRes interface{}) error {
	if err := json.Unmarshal(pData, pRes); err != nil {
		return errors.Join(ErrDeserialize, err)
	}
	return nil
}
