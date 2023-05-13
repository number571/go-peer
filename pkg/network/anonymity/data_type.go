package anonymity

import "bytes"

type iDataType byte

const (
	cIsRequest  iDataType = '>'
	cIsResponse iDataType = '<'
)

func isRequest(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return iDataType(pBytes[0]) == cIsRequest
}

func isResponse(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return iDataType(pBytes[0]) == cIsResponse
}

func wrapRequest(pBytes []byte) []byte {
	return bytes.Join([][]byte{
		{byte(cIsRequest)},
		pBytes,
	}, []byte{})
}

func wrapResponse(pBytes []byte) []byte {
	return bytes.Join([][]byte{
		{byte(cIsResponse)},
		pBytes,
	}, []byte{})
}

func unwrapBytes(pBytes []byte) []byte {
	if len(pBytes) == 0 {
		panic("length of bytes = 0")
	}
	return pBytes[1:]
}
