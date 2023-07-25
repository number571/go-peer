package flag

import "testing"

func TestFlagValue(t *testing.T) {
	argsSlice := []string{
		"--key", "123",
		"-name", "number",
		"-null", "some-value",
		"value", "571",
		"asdfg=12345",
		"-qwerty=67890",
		"--zxcvb=!@#$%",
	}

	if getFlagValueBySlice(argsSlice, "key", "1") != "123" {
		t.Error("key != 123")
		return
	}

	if getFlagValueBySlice(argsSlice, "name", "2") != "number" {
		t.Error("name != number")
		return
	}

	if getFlagValueBySlice(argsSlice, "value", "3") != "571" {
		t.Error("value != 571")
		return
	}

	if getFlagValueBySlice(argsSlice, "asdfg", "4") != "12345" {
		t.Error("asdfg != 12345")
		return
	}

	if getFlagValueBySlice(argsSlice, "qwerty", "5") != "67890" {
		t.Error("qwerty != 67890")
		return
	}

	if getFlagValueBySlice(argsSlice, "zxcvb", "6") != "!@#$%" {
		t.Error("zxcvb != !@#$%")
		return
	}

	if getFlagValueBySlice(argsSlice, "unknown", "7") != "7" {
		t.Error("unknown != 7")
		return
	}
}
