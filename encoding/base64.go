package encoding

import "encoding/base64"

// Standart encoding in package.
func Base64Encode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Standart decoding in package.
func Base64Decode(data string) []byte {
	result, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	return result
}
