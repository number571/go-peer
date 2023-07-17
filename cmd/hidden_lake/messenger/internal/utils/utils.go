package utils

import (
	_ "embed"
	"encoding/json"
	"strings"
	"unicode"
)

type sEmojis struct {
	Emojis []struct {
		Emoji     string `json:"emoji"`
		Shortname string `json:"shortname"`
	} `json:"emojis"`
}

var (
	//go:embed emoji.json
	gEmojisJSON []byte

	gEmojiReplacer *strings.Replacer
)

func init() {
	emojis := new(sEmojis)
	if err := json.Unmarshal(gEmojisJSON, emojis); err != nil {
		panic(err)
	}

	replacerList := make([]string, 0, len(emojis.Emojis))
	for _, emoji := range emojis.Emojis {
		replacerList = append(replacerList, emoji.Shortname, emoji.Emoji)
	}

	gEmojiReplacer = strings.NewReplacer(replacerList...)
}

func ReplaceTextToEmoji(s string) string {
	return gEmojiReplacer.Replace(s)
}

func HasNotWritableCharacters(str string) bool {
	for _, c := range str {
		if !unicode.IsGraphic(c) {
			return true
		}
	}
	return false
}

func GetOnlyWritableCharacters(str string) string {
	s := make([]rune, 0, len(str))
	for _, c := range str {
		if unicode.IsGraphic(c) {
			s = append(s, c)
		}
	}
	return string(s)
}
