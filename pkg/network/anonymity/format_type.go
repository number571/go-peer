package anonymity

import "bytes"

type IFormatType byte

const (
	CIsRequest  IFormatType = '>'
	CIsResponse IFormatType = '<'
)

func isRequest(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return IFormatType(pBytes[0]) == CIsRequest
}

func isResponse(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return IFormatType(pBytes[0]) == CIsResponse
}

func wrapRequest(pBytes []byte) []byte {
	return bytes.Join([][]byte{
		{byte(CIsRequest)},
		pBytes,
	}, []byte{})
}

func unwrapBytes(pBytes []byte) []byte {
	if len(pBytes) == 0 {
		panic("length of bytes = 0")
	}
	return pBytes[1:]
}

func wrapResponse(pBytes []byte) []byte {
	return bytes.Join([][]byte{
		{byte(CIsResponse)},
		pBytes,
	}, []byte{})
}
