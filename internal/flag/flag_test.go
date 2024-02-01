package flag

import "testing"

func TestPanicFlagValue(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	argsSlice := []string{
		"--key",
	}
	_ = GetFlagValue(argsSlice, "key", "_")
}

func TestBoolFlagValue(t *testing.T) {
	t.Parallel()

	argsSlice := []string{
		"--key", "123",
		"-name",
		"value", "571",
	}

	if !GetBoolFlagValue(argsSlice, "key") {
		t.Error("!key")
		return
	}

	if !GetBoolFlagValue(argsSlice, "123") {
		t.Error("!123")
		return
	}

	if !GetBoolFlagValue(argsSlice, "name") {
		t.Error("!name")
		return
	}

	if !GetBoolFlagValue(argsSlice, "value") {
		t.Error("!value")
		return
	}

	if !GetBoolFlagValue(argsSlice, "571") {
		t.Error("!571")
		return
	}

	if GetBoolFlagValue(argsSlice, "undefined") {
		t.Error("success get undefined value")
		return
	}
}

func TestFlagValue(t *testing.T) {
	t.Parallel()

	argsSlice := []string{
		"--key", "123",
		"-name", "number",
		"-null", "some-value",
		"value", "571",
		"asdfg=12345",
		"-qwerty=67890",
		"--zxcvb=!@#$%",
	}

	if GetFlagValue(argsSlice, "key", "1") != "123" {
		t.Error("key != 123")
		return
	}

	if GetFlagValue(argsSlice, "name", "2") != "number" {
		t.Error("name != number")
		return
	}

	if GetFlagValue(argsSlice, "value", "3") != "571" {
		t.Error("value != 571")
		return
	}

	if GetFlagValue(argsSlice, "asdfg", "4") != "12345" {
		t.Error("asdfg != 12345")
		return
	}

	if GetFlagValue(argsSlice, "qwerty", "5") != "67890" {
		t.Error("qwerty != 67890")
		return
	}

	if GetFlagValue(argsSlice, "zxcvb", "6") != "!@#$%" {
		t.Error("zxcvb != !@#$%")
		return
	}

	if GetFlagValue(argsSlice, "unknown", "7") != "7" {
		t.Error("unknown != 7")
		return
	}
}
