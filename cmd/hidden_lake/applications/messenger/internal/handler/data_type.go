package handler

import (
	"bytes"
	"encoding/base64"
	"html"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/internal/chars"
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

func unwrapText(pBytes []byte, pEscape bool) string {
	if len(pBytes) == 0 { // need use first isText
		panic("length of bytes = 0")
	}
	text := utils.ReplaceTextToEmoji(string(pBytes[1:]))
	if pEscape {
		return html.EscapeString(text)
	}
	return text
}

func unwrapFile(pBytes []byte, pEscape bool) (string, string) {
	if len(pBytes) == 0 { // need use first isFile
		panic("length of bytes = 0")
	}
	splited := bytes.Split(pBytes[1:], []byte{hlm_settings.CIsFile})
	if len(splited) < 2 {
		return "", ""
	}
	filename := string(splited[0])
	if chars.HasNotGraphicCharacters(filename) {
		return "", ""
	}
	fileBytes := bytes.Join(splited[1:], []byte{hlm_settings.CIsFile})
	if len(fileBytes) == 0 {
		return "", ""
	}
	base64FileBytes := base64.StdEncoding.EncodeToString(fileBytes)
	if pEscape {
		return html.EscapeString(filename), base64FileBytes
	}
	return filename, base64FileBytes
}
