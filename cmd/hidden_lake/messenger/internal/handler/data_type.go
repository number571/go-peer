package handler

import (
	"bytes"
	"encoding/base64"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
)

type iDataType byte

const (
	cIsText iDataType = 1
	cIsFile iDataType = 2
)

func isText(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return iDataType(pBytes[0]) == cIsText
}

func isFile(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return iDataType(pBytes[0]) == cIsFile
}

func wrapText(pMsg string) []byte {
	return bytes.Join([][]byte{
		{byte(cIsText)},
		[]byte(pMsg),
	}, []byte{})
}

func wrapFile(filename string, pBytes []byte) []byte {
	return bytes.Join([][]byte{
		{byte(cIsFile)},
		[]byte(filename),
		{byte(cIsFile)},
		[]byte(base64.StdEncoding.EncodeToString(pBytes)),
	}, []byte{})
}

func unwrapText(pBytes []byte) string {
	if len(pBytes) == 0 { // need use first isText
		panic("length of bytes = 0")
	}
	return string(pBytes[1:])
}

func unwrapFile(pBytes []byte) (string, string) {
	if len(pBytes) == 0 { // need use first isFile
		panic("length of bytes = 0")
	}
	splited := bytes.Split(pBytes[1:], []byte{byte(cIsFile)})
	if len(splited) < 2 {
		return "", ""
	}
	filename := string(splited[0])
	if utils.HasNotWritableCharacters(filename) {
		return "", ""
	}
	fileBytes := string(splited[1]) // in base64
	if len(fileBytes) == 0 {
		return "", ""
	}
	if utils.HasNotWritableCharacters(fileBytes) {
		return "", ""
	}
	if _, err := base64.StdEncoding.DecodeString(fileBytes); err != nil {
		return "", ""
	}
	return filename, fileBytes
}
