package encoding

import (
	"encoding/json"

	"github.com/number571/go-peer/pkg/errors"
)

func Serialize(pData interface{}, withIndent bool) []byte {
	var (
		res []byte
		err error
	)

	switch withIndent {
	case true:
		res, err = json.MarshalIndent(pData, "", "\t")
	case false:
		res, err = json.Marshal(pData)
	}

	if err != nil {
		return nil
	}
	return res
}

func Deserialize(pData []byte, pRes interface{}) error {
	if err := json.Unmarshal(pData, pRes); err != nil {
		return errors.WrapError(err, "unmarshal data")
	}
	return nil
}
