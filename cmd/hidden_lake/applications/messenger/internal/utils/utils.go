package utils

import (
	_ "embed"
	"encoding/json"
	"strings"
	"unicode"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/web"
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
	emojiSimple := new(sEmojis)
	if err := json.Unmarshal(web.GEmojiSimpleJSON, emojiSimple); err != nil {
		panic(err)
	}

	emoji := new(sEmojis)
	if err := json.Unmarshal(web.GEmojiJSON, emoji); err != nil {
		panic(err)
	}

	replacerList := make([]string, 0, len(emojiSimple.Emojis)+len(emoji.Emojis))

	for _, emoji := range emojiSimple.Emojis {
		replacerList = append(replacerList, emoji.Shortname, emoji.Emoji)
	}
	for _, emoji := range emoji.Emojis {
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
