package handler

import (
	"bytes"
	"encoding/base64"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
)

func isText(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return pBytes[0] == hlm_settings.CIsText
}

func isFile(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return pBytes[0] == hlm_settings.CIsFile
}

func wrapText(pMsg string) []byte {
	return bytes.Join([][]byte{
		{hlm_settings.CIsText},
		[]byte(pMsg),
	}, []byte{})
}

func wrapFile(filename string, pBytes []byte) []byte {
	return bytes.Join([][]byte{
		{hlm_settings.CIsFile},
		[]byte(filename),
		{hlm_settings.CIsFile},
		pBytes,
	}, []byte{})
}

func unwrapText(pBytes []byte) string {
	if len(pBytes) == 0 { // need use first isText
		panic("length of bytes = 0")
	}
	return utils.ReplaceTextToEmoji(string(pBytes[1:]))
}

func unwrapFile(pBytes []byte) (string, string) {
	if len(pBytes) == 0 { // need use first isFile
		panic("length of bytes = 0")
	}
	splited := bytes.Split(pBytes[1:], []byte{hlm_settings.CIsFile})
	if len(splited) < 2 {
		return "", ""
	}
	filename := string(splited[0])
	if utils.HasNotWritableCharacters(filename) {
		return "", ""
	}
	fileBytes := bytes.Join(splited[1:], []byte{hlm_settings.CIsFile}) // in base64
	if len(fileBytes) == 0 {
		return "", ""
	}
	return filename, base64.StdEncoding.EncodeToString(fileBytes)
}
