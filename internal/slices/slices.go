package slices

func UniqAppendToSlice(pSlice []string, pStr string) []string {
	if HasInSlice(pSlice, pStr) {
		return pSlice
	}
	return append(pSlice, pStr)
}

func DeleteFromSlice(pSlice []string, pStr string) []string {
	result := make([]string, 0, len(pSlice))
	for _, s := range pSlice {
		if s == pStr {
			continue
		}
		result = append(result, s)
	}
	return result
}

func HasInSlice(pSlice []string, pStr string) bool {
	for _, s := range pSlice {
		if pStr == s {
			return true
		}
	}
	return false
}
