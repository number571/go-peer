package encoding

import (
	yaml "gopkg.in/yaml.v2"
)

func SerializeYAML(pData interface{}) []byte {
	res, err := yaml.Marshal(pData)
	if err != nil {
		return nil
	}
	return res
}

func DeserializeYAML(pData []byte, pRes interface{}) error {
	return yaml.Unmarshal(pData, pRes)
}
