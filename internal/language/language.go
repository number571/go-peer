package language

import (
	"strings"
)

const (
	cLangENGs = "ENG"
	cLangRUSs = "RUS"
	cLangESPs = "ESP"
)

func ToILanguage(s string) (ILanguage, error) {
	switch strings.ToUpper(s) {
	case "", cLangENGs:
		return CLangENG, nil
	case cLangRUSs:
		return CLangRUS, nil
	case cLangESPs:
		return CLangESP, nil
	default:
		return 0, ErrUnknownLanguage
	}
}

func FromILanguage(pLang ILanguage) string {
	switch pLang {
	case CLangENG:
		return cLangENGs
	case CLangRUS:
		return cLangRUSs
	case CLangESP:
		return cLangESPs
	default:
		panic("unknown language")
	}
}
