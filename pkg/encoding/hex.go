package encoding

import "encoding/hex"

func HexEncode(pData []byte) string {
	return hex.EncodeToString(pData)
}

func HexDecode(pData string) []byte {
	result, err := hex.DecodeString(pData)
	if err != nil {
		return nil
	}
	return result
}
