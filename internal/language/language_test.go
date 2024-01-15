package language

import "testing"

func TestPanicFromLanguage(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = FromILanguage(111)
}

func TestToLanguage(t *testing.T) {
	lang, err := ToILanguage("ENG")
	if err != nil {
		t.Error(err)
		return
	}
	if lang != CLangENG {
		t.Error("got invalid ENG")
		return
	}

	lang, err = ToILanguage("RUS")
	if err != nil {
		t.Error(err)
		return
	}
	if lang != CLangRUS {
		t.Error("got invalid RUS")
		return
	}

	lang, err = ToILanguage("ESP")
	if err != nil {
		t.Error(err)
		return
	}
	if lang != CLangESP {
		t.Error("got invalid ESP")
		return
	}

	if _, err := ToILanguage("???"); err == nil {
		t.Error("success unknown type to language")
		return
	}
}

func TestFromLanguage(t *testing.T) {
	if FromILanguage(CLangENG) != "ENG" {
		t.Error("got invalid ENG")
		return
	}
	if FromILanguage(CLangRUS) != "RUS" {
		t.Error("got invalid RUS")
		return
	}
	if FromILanguage(CLangESP) != "ESP" {
		t.Error("got invalid ESP")
		return
	}
}
