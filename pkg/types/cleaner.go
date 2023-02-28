package types

// returns last error from slice
func CleanAll(cs []ICleaner) error {
	var lastErr error
	for _, c := range cs {
		if err := c.Clean(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
