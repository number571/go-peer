package msgconv

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
)

func FromStringToBytes(pMsg string) []byte {
	splited := strings.Split(pMsg, message.CSeparator)
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
		[]byte(message.CSeparator),
	)
}

func FromBytesToString(pMsg []byte) string {
	splited := bytes.Split(pMsg, []byte(message.CSeparator))
	if len(splited) < 2 {
		return ""
	}
	encBytes := encoding.HexEncode(bytes.Join(splited[1:], []byte(message.CSeparator)))
	return strings.Join(
		[]string{
			removeInvisibleChars(string(splited[0])),
			encBytes,
		},
		message.CSeparator,
	)
}

func removeInvisibleChars(pS string) string {
	return strings.TrimFunc(pS, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})
}
