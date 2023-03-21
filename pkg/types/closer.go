package types

// returns last error from slice
func CloseAll(pClosers []ICloser) error {
	var lastErr error
	for _, c := range pClosers {
		if err := c.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
