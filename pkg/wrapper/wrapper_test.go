package wrapper

import "testing"

func TestWrapper(t *testing.T) {
	t.Parallel()

	wr := NewWrapper()
	val := "test_1"

	if val != wr.Set(val).Get().(string) {
		t.Error("got value not equal original")
	}

	newVal := "test_2"
	if newVal != wr.Set(newVal).Get().(string) {
		t.Error("got new value not equal new original")
	}
}
