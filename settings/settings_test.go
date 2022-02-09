package settings

import "testing"

func TestSettings(t *testing.T) {
	s := NewSettings()
	v := s.Get(SizeWork)

	if s.Get(SizeWork) != v {
		t.Errorf("value is not determined")
	}

	v += 1
	s.Set(SizeWork, v)

	if s.Get(SizeWork) != v {
		t.Errorf("value is not saved")
	}
}
