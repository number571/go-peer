package settings

import "testing"

func TestSettings(t *testing.T) {
	s := NewSettings()
	v := s.Get(CSizeWork)

	if s.Get(CSizeWork) != v {
		t.Errorf("value is not determined")
	}

	v += 1
	s.Set(CSizeWork, v)

	if s.Get(CSizeWork) != v {
		t.Errorf("value is not saved")
	}
}
