package settings

import "testing"

func TestSettings(t *testing.T) {
	v := GSettings.Get(CSizePayl).(uint64)

	if GSettings.Get(CSizePayl).(uint64) != v {
		t.Errorf("value is not determined")
	}

	v += 1
	GSettings.Set(CSizePayl, v)

	if GSettings.Get(CSizePayl).(uint64) != v {
		t.Errorf("value is not saved")
	}
}
