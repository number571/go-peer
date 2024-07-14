package chars

import "unicode"

func HasNotGraphicCharacters(pS string) bool {
	for _, c := range pS {
		if !unicode.IsGraphic(c) {
			return true
		}
	}
	return false
}
