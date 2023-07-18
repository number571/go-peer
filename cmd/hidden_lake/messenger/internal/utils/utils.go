package utils

import (
	_ "embed"
	"encoding/json"
	"strings"
	"unicode"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
)

type sEmojis struct {
	Emojis []struct {
		Emoji     string `json:"emoji"`
		Shortname string `json:"shortname"`
	} `json:"emojis"`
}

var (
	gEmojiReplacer *strings.Replacer
)

func init() {
	emojis := new(sEmojis)
	if err := json.Unmarshal(web.GEmojisJSON, emojis); err != nil {
		panic(err)
	}

	replacerList := make([]string, 0, len(emojis.Emojis))
	for _, emoji := range emojis.Emojis {
		replacerList = append(replacerList, emoji.Shortname, emoji.Emoji)
	}

	gEmojiReplacer = strings.NewReplacer(replacerList...)
}

func ReplaceTextToEmoji(pS string) string {
	return gEmojiReplacer.Replace(pS)
}

func HasNotWritableCharacters(pS string) bool {
	for _, c := range pS {
		if !unicode.IsGraphic(c) {
			return true
		}
	}
	return false
}

func GetOnlyWritableCharacters(pS string) string {
	s := make([]rune, 0, len(pS))
	for _, c := range pS {
		if unicode.IsGraphic(c) {
			s = append(s, c)
		}
	}
	return string(s)
}
