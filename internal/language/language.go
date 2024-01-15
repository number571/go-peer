package language

import (
	"errors"
	"strings"
)

func ToILanguage(s string) (ILanguage, error) {
	switch strings.ToUpper(s) {
	case "", "ENG":
		return CLangENG, nil
	case "RUS":
		return CLangRUS, nil
	case "ESP":
		return CLangESP, nil
	default:
		return 0, errors.New("unknown language")
	}
}

func FromILanguage(pLang ILanguage) string {
	switch pLang {
	case CLangENG:
		return "ENG"
	case CLangRUS:
		return "RUS"
	case CLangESP:
		return "ESP"
	default:
		panic("unknown language")
	}
}
