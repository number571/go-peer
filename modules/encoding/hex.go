package encoding

import "encoding/hex"

func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

func HexDecode(data string) []byte {
	result, err := hex.DecodeString(data)
	if err != nil {
		return nil
	}
	return result
}
