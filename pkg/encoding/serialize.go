package encoding

import "encoding/json"

func Serialize(pData interface{}) []byte {
	res, err := json.Marshal(pData)
	if err != nil {
		return nil
	}
	return res
}

func Deserialize(pData []byte, pRes interface{}) error {
	return json.Unmarshal(pData, pRes)
}
