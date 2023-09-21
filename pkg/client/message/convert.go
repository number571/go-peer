package message

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/number571/go-peer/pkg/encoding"
)

func FromStringToBytes(pMsg string) []byte {
	splited := strings.Split(pMsg, cSeparator)
	if len(splited) != 2 {
		return nil
	}
	decBytes := encoding.HexDecode(strings.TrimSpace(splited[1]))
	if decBytes == nil {
		return nil
	}
	return bytes.Join(
		[][]byte{
			[]byte(removeInvisibleChars(splited[0])),
			decBytes,
		},
		[]byte(cSeparator),
	)
}

func FromBytesToString(pMsg []byte) string {
	splited := bytes.Split(pMsg, []byte(cSeparator))
	if len(splited) < 2 {
		return ""
	}
	encBytes := encoding.HexEncode(bytes.Join(splited[1:], []byte(cSeparator)))
	return strings.Join(
		[]string{
			removeInvisibleChars(string(splited[0])),
			encBytes,
		},
		cSeparator,
	)
}

func removeInvisibleChars(pS string) string {
	return strings.TrimFunc(pS, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})
}
