package stringtools

import "testing"

func TestDeleteFromSlice(t *testing.T) {
	t.Parallel()

	slice := []string{"a", "b", "c", "d"}
	elem := "c"
	result := DeleteFromSlice(slice, elem)
	for _, v := range result {
		if v == elem {
			t.Error("found deleted element")
			return
		}
	}
}

func TestUniqAppendToSlice(t *testing.T) {
	t.Parallel()

	slice := []string{"a", "b", "c", "d"}
	elem1 := "c"
	elem2 := "e"

	slice = UniqAppendToSlice(slice, elem1)
	slice = UniqAppendToSlice(slice, elem2)

	count := 0
	newElemFound := false
	for _, v := range slice {
		if v == elem1 {
			count++
			if count == 2 {
				t.Error("found not uniq element")
				return
			}
		}
		if v == elem2 {
			newElemFound = true
			break
		}
	}

	if !newElemFound {
		t.Error("error append new element")
		return
	}
}
