package encoding

import "encoding/json"

func Serialize(data interface{}) []byte {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil
	}
	return res
}

func Deserialize(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
