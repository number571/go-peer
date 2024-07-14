package chars

import "testing"

func TestHasNotGraphicCharacters(t *testing.T) {
	t.Parallel()

	if HasNotGraphicCharacters("hello, world!") {
		t.Error("message contains only graphic chars")
		return
	}

	if !HasNotGraphicCharacters("hello,\nworld!") {
		t.Error("message contains not graphic chars")
		return
	}
}
